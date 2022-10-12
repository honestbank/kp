package kp_test

import (
	"context"
	"errors"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/honestbank/kp"
	"github.com/stretchr/testify/assert"
)

const retiresCount = 2

var kafkaConfig = kp.KafkaConfig{
	KafkaBootstrapServers: []string{"localhost:9092"},
}

var data = SafeCounter{v: make(map[string]int)}
var producer = kp.NewProducer(kafkaConfig)

type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

// Inc increments the counter for the given key.
func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[key]++
	c.mu.Unlock()
}

// Value returns the current value of the counter for the given key.
func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mu.Unlock()

	return c.v[key]
}

func TestNewKafkaProcessor(t *testing.T) {
	t.Run("test new kafka processor", func(t *testing.T) {
		a := assert.New(t)

		processor := kp.NewKafkaProcessor("test", "retry-test", "dead-test", 10, "test", kafkaConfig, 0)
		a.NotNil(processor)
	})

	t.Run("test process", func(t *testing.T) {
		ctx, ctxWithCancel := context.WithCancel(context.Background())
		processor := CreateKPProcessor()
		go produceMockMessage("success")
		go produceMockMessage("success")
		go produceMockMessage("success")
		go sendContextDoneSignal(10, ctxWithCancel)
		_ = processor.Start(ctx)
		assert.Equal(t, 3, data.Value("success"))
	})

	t.Run("test process fail", func(t *testing.T) {
		processor := CreateKPProcessor()
		go produceMockMessage("fail")
		go sendShutdownSignal(10)
		_ = processor.Start(context.Background()) //cannot handle error due to not implement channel in this version(waiting for v2)
		assert.Equal(t, retiresCount, data.Value("fail"))
	})
}

func sendShutdownSignal(duration time.Duration) {
	time.Sleep(duration * time.Second)
	_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
}

func sendContextDoneSignal(duration time.Duration, ctx context.CancelFunc) {
	time.Sleep(duration * time.Second)
	ctx()
}

func ReceiveMessage(ctx context.Context, key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error {
	if message == "fail" {
		data.Inc("fail")

		return errors.New("failed message from ReceivedMessage")
	}
	data.Inc("success")

	return nil
}

func CreateKPProcessor() kp.KafkaProcessor {
	data = SafeCounter{v: make(map[string]int)}
	p := kp.NewKafkaProcessor(
		"test-topic", "retry-test", "dead-test",
		retiresCount,
		"test-application-name",
		kafkaConfig,
		time.Second*1)
	p.Process(ReceiveMessage)

	return p
}

func produceMockMessage(msg string) {
	time.Sleep(time.Second * 8) // waiting consumer ready
	_ = producer.ProduceMessage(context.Background(), "test-topic", "1", msg)
}
