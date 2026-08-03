package kafkaheaders_test

import (
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/kafkaheaders"
)

func TestGet(t *testing.T) {
	t.Run("returns nil when header is not found", func(t *testing.T) {
		assert.Nil(t, kafkaheaders.Get("retry-count", &kafka.Message{}))
	})
	t.Run("returns nil if the header value is nil", func(t *testing.T) {
		assert.Nil(t, kafkaheaders.Get("retry-count", &kafka.Message{Headers: []kafka.Header{{Key: "something", Value: []byte("s")}, {Key: "retry-count", Value: nil}}}))
	})
	t.Run("returns value if the header value is nil", func(t *testing.T) {
		assert.Equal(t, "count", *kafkaheaders.Get("retry-count", &kafka.Message{Headers: []kafka.Header{{Key: "something", Value: []byte("s")}, {Key: "retry-count", Value: []byte("count")}}}))
	})
}
