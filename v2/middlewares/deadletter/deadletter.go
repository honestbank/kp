package deadletter

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/honestbank/kp/v2/internal/retrycounter"
	"github.com/honestbank/kp/v2/middlewares"
)

type deadletter struct {
	producer        Producer
	onProduceErrors func(err error)
	threshold       int
}
type Producer interface {
	ProduceRaw(message *kafka.Message) error
}

func (r deadletter) Process(ctx context.Context, item *kafka.Message, next func(ctx context.Context, item *kafka.Message) error) error {
	err := next(ctx, item)
	if err == nil {
		return nil
	}
	if retrycounter.GetCount(item) < r.threshold {
		return err
	}
	err = r.producer.ProduceRaw(&kafka.Message{Value: item.Value, Key: item.Key, Headers: item.Headers, Timestamp: item.Timestamp, TimestampType: item.TimestampType, Opaque: item.Opaque})
	if err != nil {
		r.onProduceErrors(err)
	}

	return nil
}

func NewDeadletterMiddleware(producer Producer, threshold int, onProduceErrors func(error)) middlewares.KPMiddleware[*kafka.Message] {
	return deadletter{onProduceErrors: onProduceErrors, producer: producer, threshold: threshold}
}
