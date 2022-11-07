package kafkaheaders_test

import (
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/internal/kafkaheaders"
)

func TestGet(t *testing.T) {
	t.Run("returns nil when header is not found", func(t *testing.T) {
		assert.Nil(t, kafkaheaders.Get("retry-count", &kafka.Message{}))
	})
	t.Run("returns nil if the header value is nil", func(t *testing.T) {
		assert.Nil(t, kafkaheaders.Get("retry-count", &kafka.Message{Headers: []kafka.Header{{"something", []byte("s")}, {"retry-count", nil}}}))
	})
	t.Run("returns value if the header value is nil", func(t *testing.T) {
		assert.Equal(t, "count", *kafkaheaders.Get("retry-count", &kafka.Message{Headers: []kafka.Header{{"something", []byte("s")}, {"retry-count", []byte("count")}}}))
	})
}
