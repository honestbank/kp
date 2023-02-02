package deadletter_test

import (
	"context"
	"errors"
	"testing"

	"github.com/honestbank/kp/v2/internal/retrycounter"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/middlewares/deadletter"
	"github.com/honestbank/kp/v2/producer"
)

type producerMock struct {
	produceRaw func(message *kafka.Message) error
}

func (r producerMock) Flush() error {
	return nil
}

func (r producerMock) Produce(context context.Context, message any) error {
	return nil
}

func (r producerMock) ProduceRaw(message *kafka.Message) error {
	return r.produceRaw(message)
}

func newProducer(cb func(item *kafka.Message) error) producer.Producer[any] {
	return producerMock{produceRaw: cb}
}

func TestRetry_Process(t *testing.T) {
	t.Run("if message processing succeeds, it returns nil without retrying", func(t *testing.T) {
		middleware := deadletter.NewDeadletterMiddleware(nil, 2, nil)
		assert.NotPanics(t, func() {
			err := middleware.Process(context.Background(), nil, func(ctx context.Context, item *kafka.Message) error {
				return nil
			})
			assert.NoError(t, err)
		})
	})

	t.Run("if the retry count is less than threshold it simply returns error", func(t *testing.T) {
		middleware := deadletter.NewDeadletterMiddleware(newProducer(func(item *kafka.Message) error {
			return nil
		}), 2, func(err error) {
			t.FailNow()
		})
		msg := &kafka.Message{}
		retrycounter.SetCount(msg, 0)
		err := middleware.Process(context.Background(), msg, func(ctx context.Context, item *kafka.Message) error {
			return errors.New("random error")
		})
		assert.Error(t, err)
		retrycounter.SetCount(msg, 1)
		err = middleware.Process(context.Background(), msg, func(ctx context.Context, item *kafka.Message) error {
			return errors.New("random error")
		})
		assert.Error(t, err)
		retrycounter.SetCount(msg, 2)
		err = middleware.Process(context.Background(), msg, func(ctx context.Context, item *kafka.Message) error {
			return errors.New("random error")
		})
		assert.NoError(t, err)
	})

	t.Run("if producing returns error, we get a callback", func(t *testing.T) {
		called := false
		middleware := deadletter.NewDeadletterMiddleware(newProducer(func(item *kafka.Message) error {
			return errors.New("random error")
		}), 1, func(err error) {
			called = true
		})
		msg := &kafka.Message{}
		retrycounter.SetCount(msg, 2)
		middleware.Process(context.Background(), msg, func(ctx context.Context, item *kafka.Message) error {
			return errors.New("random error")
		})
		assert.True(t, called)
	})
}
