package retry_count_test

import (
	"context"
	"testing"

	"github.com/honestbank/kp/v2/middlewares/retry_count"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/internal/retrycounter"
)

func TestRetryCount(t *testing.T) {
	t.Run("sets retry count if there's no value at all", func(t *testing.T) {
		_ = retry_count.NewRetryCountMiddleware().Process(context.Background(), &kafka.Message{}, func(ctx context.Context, item *kafka.Message) error {
			assert.Equal(t, 0, retry_count.FromContext(ctx))
			return nil
		})
	})
	t.Run("returns correct retry count", func(t *testing.T) {
		message := &kafka.Message{}
		retrycounter.SetCount(message, 50)
		_ = retry_count.NewRetryCountMiddleware().Process(context.Background(), message, func(ctx context.Context, item *kafka.Message) error {
			assert.Equal(t, 50, retry_count.FromContext(ctx))
			return nil
		})
	})
}

func TestRetryCountFromContext(t *testing.T) {
	t.Run("returns 0 if there's no value", func(t *testing.T) {
		assert.Equal(t, 0, retry_count.FromContext(context.Background()))
	})
}
