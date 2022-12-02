package kafkaheaders_test

import (
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/kafkaheaders"
)

func TestSet(t *testing.T) {
	t.Run("can set when there's no existing header", func(t *testing.T) {
		msg := kafka.Message{}
		kafkaheaders.Set(&msg, "key", "value")
		assert.Equal(t, "value", *kafkaheaders.Get("key", &msg))
	})
	t.Run("can set multiple keys and overwrite existing header", func(t *testing.T) {
		msg := kafka.Message{}
		kafkaheaders.Set(&msg, "key", "value")
		kafkaheaders.Set(&msg, "key2", "value2")
		assert.Equal(t, "value", *kafkaheaders.Get("key", &msg))
		assert.Equal(t, "value2", *kafkaheaders.Get("key2", &msg))
		kafkaheaders.Set(&msg, "key", "new-value")
		kafkaheaders.Set(&msg, "key2", "new-value2")
		assert.Equal(t, "new-value", *kafkaheaders.Get("key", &msg))
		assert.Equal(t, "new-value2", *kafkaheaders.Get("key2", &msg))
	})
	t.Run("can work with pre-existing header", func(t *testing.T) {
		msg := kafka.Message{
			Headers: []kafka.Header{{Key: "key", Value: []byte("value")}},
		}
		kafkaheaders.Set(&msg, "key", "value")
		kafkaheaders.Set(&msg, "key2", "value2")
		assert.Equal(t, "value", *kafkaheaders.Get("key", &msg))
		assert.Equal(t, "value2", *kafkaheaders.Get("key2", &msg))
		kafkaheaders.Set(&msg, "key", "new-value")
		kafkaheaders.Set(&msg, "key2", "new-value2")
		assert.Equal(t, "new-value", *kafkaheaders.Get("key", &msg))
		assert.Equal(t, "new-value2", *kafkaheaders.Get("key2", &msg))
	})
}
