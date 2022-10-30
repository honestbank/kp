package retrycounter_test

import (
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/internal/retrycounter"
)

func TestGetCount(t *testing.T) {
	t.Run("returns 0 when there's no count in header", func(t *testing.T) {
		assert.Equal(t, 0, retrycounter.GetCount(&kafka.Message{}))
	})
	t.Run("returns correct count when there's count in header", func(t *testing.T) {
		assert.Equal(t, 5, retrycounter.GetCount(&kafka.Message{Headers: []kafka.Header{{Key: "x-retry-count", Value: []byte("5")}}}))
	})
}

func TestSetCount(t *testing.T) {
	t.Run("Can set when there's no retry count in header", func(t *testing.T) {
		message := &kafka.Message{}
		retrycounter.SetCount(message, 10)
		assert.Equal(t, 10, retrycounter.GetCount(message))
	})

	t.Run("returns correct count when there's count in header", func(t *testing.T) {
		assert.Equal(t, 5, retrycounter.GetCount(&kafka.Message{Headers: []kafka.Header{{Key: "x-retry-count", Value: []byte("5")}}}))
	})
}
