//go:build integration_test

package consumer_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/internal/consumer"
	"github.com/honestbank/kp/v2/producer"
)

type MyMsg struct {
	Time string
}

func TestNew(t *testing.T) {
	t.Setenv("KP_SCHEMA_REGISTRY_ENDPOINT", "http://localhost:8081")
	t.Setenv("KP_KAFKA_BOOTSTRAP_SERVERS", "localhost")
	t.Run("can read from kafka", func(t *testing.T) {
		c, err := consumer.New([]string{"consumer-integration-topic-1"}, "consumer-group-1")
		assert.NoError(t, err)
		p1, err := producer.New[MyMsg]("consumer-integration-topic-1")
		assert.NoError(t, err)
		shouldContinue, numberOfMessage := true, 0
		go func() {
			for shouldContinue {
				message := c.GetMessage()
				if message == nil {
					continue
				}
				c.Commit(message)
				numberOfMessage++
			}
		}()
		err = p1.Produce(context.Background(), MyMsg{Time: time.Now().Format(time.RFC3339Nano)})
		assert.NoError(t, err)
		err = p1.Produce(context.Background(), MyMsg{Time: time.Now().Format(time.RFC3339Nano)})
		assert.NoError(t, err)
		err = p1.Produce(context.Background(), MyMsg{Time: time.Now().Format(time.RFC3339Nano)})
		assert.NoError(t, err)
		time.Sleep(time.Millisecond * 500)
		p1.Flush()
		time.Sleep(time.Millisecond * 1000)
		shouldContinue = false
		time.Sleep(time.Millisecond * 500)
		assert.Equal(t, 3, numberOfMessage)
	})
	t.Run("can read from multiple topics", func(t *testing.T) {
		c, err := consumer.New([]string{"consumer-integration-topic-2", "consumer-integration-topic-3"}, "consumer-group-2")
		assert.NoError(t, err)
		p1, err := producer.New[MyMsg]("consumer-integration-topic-2")
		assert.NoError(t, err)
		p2, err := producer.New[MyMsg]("consumer-integration-topic-3")
		assert.NoError(t, err)
		shouldContinue, numberOfMessage := true, 0
		go func() {
			for shouldContinue {
				message := c.GetMessage()
				if message != nil {
					c.Commit(message)
					numberOfMessage++
				}
			}
		}()
		err = p1.Produce(context.Background(), MyMsg{Time: time.Now().Format(time.RFC3339Nano)})
		assert.NoError(t, err)
		err = p2.Produce(context.Background(), MyMsg{Time: time.Now().Format(time.RFC3339Nano)})
		assert.NoError(t, err)
		err = p2.Produce(context.Background(), MyMsg{Time: time.Now().Format(time.RFC3339Nano)})
		assert.NoError(t, err)
		p1.Flush()
		p2.Flush()
		time.Sleep(time.Millisecond * 500)
		shouldContinue = false
		time.Sleep(time.Millisecond * 500)
		assert.Equal(t, 3, numberOfMessage)
	})
	t.Run("returns error if config is invalid", func(t *testing.T) {
		t.Setenv("KP_SCHEMA_REGISTRY_ENDPOINT", "")
		t.Setenv("KP_KAFKA_BOOTSTRAP_SERVERS", "")
		c, err := consumer.New([]string{}, "")
		assert.Error(t, err)
		assert.Nil(t, c)
	})
	t.Run("returns error if there's no topic", func(t *testing.T) {
		c, err := consumer.New([]string{}, "")
		assert.Error(t, err)
		assert.Nil(t, c)
	})
}
