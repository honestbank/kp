package consumer_test

import (
	"context"
	"errors"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/internal/serialization"
	"github.com/honestbank/kp/v2/middlewares/consumer"
)

type mockConsumer struct {
	count      int
	returnsNil bool
}
type Payload struct {
	Count int
}

func (m *mockConsumer) GetMessage() *kafka.Message {
	if m.returnsNil {
		return nil
	}
	m.count = m.count + 1
	encode, err := serialization.Encode(Payload{Count: m.count}, 1)
	if err != nil {
		return nil
	}
	return &kafka.Message{
		Value: encode,
	}
}

func (m *mockConsumer) Commit(message *kafka.Message) error {
	return nil
}
func TestConsumerMiddleware_Process(t *testing.T) {
	t.Run("if the message is not nil, it doesn't call GetMessage", func(t *testing.T) {
		c := &mockConsumer{}
		middleware := consumer.NewConsumerMiddleware(c)
		msg := kafka.Message{Value: []byte("original-value")}
		err := middleware.Process(context.Background(), &msg, func(ctx context.Context, item *kafka.Message) error {
			assert.Equal(t, string(item.Value), string(msg.Value))
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 0, c.count)
	})
	t.Run("if message is nil, it uses the message from consumer", func(t *testing.T) {
		c := &mockConsumer{}
		middleware := consumer.NewConsumerMiddleware(c)
		err := middleware.Process(context.Background(), nil, func(ctx context.Context, item *kafka.Message) error {
			decode, err := serialization.Decode[Payload](item.Value)
			if err != nil {
				return err
			}
			assert.Equal(t, 1, decode.Count)
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, c.count)
	})

	t.Run("calls next and returns what next returns", func(t *testing.T) {
		c := &mockConsumer{}
		middleware := consumer.NewConsumerMiddleware(c)
		customErr := errors.New("some error")
		err := middleware.Process(context.Background(), nil, func(ctx context.Context, item *kafka.Message) error {
			return customErr
		})
		assert.ErrorIs(t, err, customErr)
	})
	t.Run("if consumer returns nil, it returns nil without calling next", func(t *testing.T) {
		c := &mockConsumer{returnsNil: true}
		middleware := consumer.NewConsumerMiddleware(c)
		err := middleware.Process(context.Background(), nil, func(ctx context.Context, item *kafka.Message) error {
			assert.FailNow(t, "next shouldn't have been called")

			return nil
		})
		assert.NoError(t, err)
	})
}
