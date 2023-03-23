package retry_test

import (
	"context"
	"errors"
	"github.com/honestbank/kp/v2/internal/middleware"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/internal/retrycounter"
	"github.com/honestbank/kp/v2/middlewares/retry"
	"github.com/honestbank/kp/v2/producer"
)

type producerMock struct {
	produceRaw func(message *kafka.Message) error
}

func (r producerMock) Flush() error {
	return nil
}

func (r producerMock) Produce(context context.Context, message *kafka.Message) error {
	return nil
}

func (r producerMock) SetMiddlewares([]middleware.Middleware[*kafka.Message, error]) {
}

func newProducer(cb func(item *kafka.Message) error) producer.UntypedProducer {
	return producerMock{produceRaw: cb}
}

func TestRetry_Process(t *testing.T) {
	t.Run("if message processing succeeds, it returns nil without retrying", func(t *testing.T) {
		middleware := retry.NewRetryMiddleware(nil, nil)
		assert.NotPanics(t, func() {
			err := middleware.Process(context.Background(), nil, func(ctx context.Context, item *kafka.Message) error {
				return nil
			})
			assert.NoError(t, err)
		})
	})

	t.Run("producing message works and increments the retry count", func(t *testing.T) {
		middleware := retry.NewRetryMiddleware(newProducer(func(item *kafka.Message) error {
			return nil
		}), func(err error) {
			t.FailNow()
		})
		msg := &kafka.Message{}
		middleware.Process(context.Background(), msg, func(ctx context.Context, item *kafka.Message) error {
			assert.Equal(t, 0, retrycounter.GetCount(item))
			return errors.New("random error")
		})
		middleware.Process(context.Background(), msg, func(ctx context.Context, item *kafka.Message) error {
			assert.Equal(t, 1, retrycounter.GetCount(item))
			return errors.New("random error")
		})
		middleware.Process(context.Background(), msg, func(ctx context.Context, item *kafka.Message) error {
			assert.Equal(t, 2, retrycounter.GetCount(item))
			return errors.New("random error")
		})
	})

	t.Run("if producing fails while retrying, we get a callback", func(t *testing.T) {
		called := false
		middleware := retry.NewRetryMiddleware(newProducer(func(item *kafka.Message) error {
			return errors.New("some error")
		}), func(err error) {
			called = true
		})
		msg := &kafka.Message{}
		middleware.Process(context.Background(), msg, func(ctx context.Context, item *kafka.Message) error {
			assert.Equal(t, 0, retrycounter.GetCount(item))
			return errors.New("random error")
		})
		assert.True(t, called)
	})
}
