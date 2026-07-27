//go:build integration_test

// End-to-end test proving the automatic consumer Close() wired into kp.Run()
// actually leaves the Kafka group on the wire.
//
// Runs in CI under the integration-tests job (real broker on localhost:9092).
// Locally:
//
//	go test -tags integration_test -run TestAutoClose -v ./...
//
// Set KP_DEBUG_CGRP=1 to enable librdkafka's cgrp debug, which prints the
// LeaveGroup request to stderr as direct wire-level evidence:
//
//	KP_DEBUG_CGRP=1 go test -tags integration_test -run TestAutoClose -v ./... 2>&1 | grep -i leavegroup
package v2_test

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v2 "github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/config"
	kpconsumer "github.com/honestbank/kp/v2/consumer"
	consumermw "github.com/honestbank/kp/v2/middlewares/consumer"
)

// autoCloseBootstrap is the broker the test talks to. CI publishes the broker on
// localhost:9092; override with KP_BOOTSTRAP for a broker on another port.
var autoCloseBootstrap = func() string {
	if v := os.Getenv("KP_BOOTSTRAP"); v != "" {
		return v
	}
	return "localhost:9092"
}()

// TestAutoClose covers the automatic consumer Close() that kp.Run() performs when
// its processing loop stops. Each subtest names the auto-close behaviour it proves.
func TestAutoClose(t *testing.T) {
	t.Run("auto-close leaves group so a replacement consumer takes over promptly", func(t *testing.T) {
		suffix := fmt.Sprintf("%d", time.Now().UnixNano())
		topic := "kp-autoclose-" + suffix
		group := "kp-autoclose-grp-" + suffix

		createAutoCloseTopic(t, topic, 1)

		// Consumer A joins the group, consumes one message, then is stopped.
		// Stopping ends Run's loop, which now calls chain.Close() -> consumer.Close()
		// -> librdkafka LeaveGroup. Enable cgrp debug (opt-in) to see it on stderr.
		cfgA := config.Kafka{BootstrapServers: autoCloseBootstrap, ConsumerGroupName: group}
		if os.Getenv("KP_DEBUG_CGRP") != "" {
			debug := "cgrp"
			cfgA.Debug = &debug
		}
		consumerA, err := kpconsumer.New([]string{topic}, cfgA.WithDefaults())
		require.NoError(t, err)

		procA := v2.New[kafka.Message]()
		var countA int64
		doneA := make(chan error, 1)
		go func() {
			doneA <- procA.
				AddMiddleware(consumermw.NewConsumerMiddleware(consumerA)).
				Run(countingHandler(&countA))
		}()

		produceRaw(t, topic, "m1")
		require.True(t, waitForCond(func() bool { return atomic.LoadInt64(&countA) >= 1 }, 40*time.Second),
			"consumer A never joined the group / received m1")
		t.Logf("consumer A joined and consumed m1")

		// Trigger graceful shutdown. Run must return (the window in which auto-close
		// + LeaveGroup happen); a hang here means auto-close blocked.
		closeStart := time.Now()
		procA.Stop()
		select {
		case runErr := <-doneA:
			require.NoError(t, runErr, "Run should return nil after graceful stop")
			t.Logf("Run returned %.2fs after Stop() (auto-close + LeaveGroup in this window)", time.Since(closeStart).Seconds())
		case <-time.After(30 * time.Second):
			t.Fatal("Run did not return within 30s after Stop() — auto-close likely blocked")
		}

		// Consumer B, same group. Because A left the group cleanly, B is assigned the
		// partition and receives a fresh message — correctness is the assertion; the
		// elapsed time is logged as evidence, not gated (CI runners vary).
		cfgB := config.Kafka{BootstrapServers: autoCloseBootstrap, ConsumerGroupName: group}
		consumerB, err := kpconsumer.New([]string{topic}, cfgB.WithDefaults())
		require.NoError(t, err)

		procB := v2.New[kafka.Message]()
		var countB int64
		doneB := make(chan error, 1)
		handoffStart := time.Now()
		go func() {
			doneB <- procB.
				AddMiddleware(consumermw.NewConsumerMiddleware(consumerB)).
				Run(countingHandler(&countB))
		}()

		produceRaw(t, topic, "m2")
		require.True(t, waitForCond(func() bool { return atomic.LoadInt64(&countB) >= 1 }, 40*time.Second),
			"consumer B never took over / received m2 — group was not left cleanly")
		t.Logf("consumer B took over and consumed m2 in %.2fs", time.Since(handoffStart).Seconds())

		procB.Stop()
		<-doneB
	})

	t.Run("auto-close is idempotent so an existing defer Close() stays safe", func(t *testing.T) {
		suffix := fmt.Sprintf("%d", time.Now().UnixNano())
		topic := "kp-autoclose-idem-" + suffix
		group := "kp-autoclose-idem-grp-" + suffix
		createAutoCloseTopic(t, topic, 1)

		cfg := config.Kafka{BootstrapServers: autoCloseBootstrap, ConsumerGroupName: group}
		c, err := kpconsumer.New([]string{topic}, cfg.WithDefaults())
		require.NoError(t, err)

		c.GetMessage() // poll once so the handle is fully live / joined

		require.NotPanics(t, func() {
			assert.NoError(t, c.Close(), "first Close")
			assert.NoError(t, c.Close(), "second Close (idempotent)")
		})
	})
}

// --- helpers ---

func countingHandler(counter *int64) func(ctx context.Context, msg *kafka.Message) error {
	return func(ctx context.Context, msg *kafka.Message) error {
		atomic.AddInt64(counter, 1)
		return nil
	}
}

func waitForCond(cond func() bool, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return cond()
}

func createAutoCloseTopic(t *testing.T, topic string, partitions int) {
	t.Helper()
	admin, err := kafka.NewAdminClient(config.GetKafkaConfig(config.Kafka{BootstrapServers: autoCloseBootstrap}))
	require.NoError(t, err)
	defer admin.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, err = admin.CreateTopics(ctx, []kafka.TopicSpecification{
		{Topic: topic, NumPartitions: partitions, ReplicationFactor: 1},
	})
	require.NoError(t, err)
}

func produceRaw(t *testing.T, topic, value string) {
	t.Helper()
	p, err := kafka.NewProducer(config.GetKafkaConfig(config.Kafka{BootstrapServers: autoCloseBootstrap}))
	require.NoError(t, err)
	defer p.Close()

	deliv := make(chan kafka.Event, 1)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(value),
	}, deliv)
	require.NoError(t, err)

	select {
	case e := <-deliv:
		m := e.(*kafka.Message)
		require.NoError(t, m.TopicPartition.Error)
	case <-time.After(15 * time.Second):
		t.Fatalf("produce of %q timed out", value)
	}
}
