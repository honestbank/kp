package retry

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/honestbank/kp/v2/internal/retrycounter"
	"github.com/honestbank/kp/v2/middlewares"
)

type retry struct {
	producer        Producer
	onProduceErrors func(err error)
}

type Producer interface {
	ProduceRaw(item *kafka.Message) error
}

func (r retry) Process(ctx context.Context, item *kafka.Message, next func(ctx context.Context, item *kafka.Message) error) error {
	err := next(ctx, item)
	if err == nil {
		return nil
	}
	retrycounter.SetCount(item, retrycounter.GetCount(item)+1)
	err = r.producer.ProduceRaw(&kafka.Message{Value: item.Value, Key: item.Key, Headers: item.Headers, Timestamp: item.Timestamp, TimestampType: item.TimestampType, Opaque: item.Opaque})
	if err != nil {
		r.onProduceErrors(err)
	}

	return nil
}

func NewRetryMiddleware(producer Producer, onProduceErrors func(error)) middlewares.KPMiddleware[*kafka.Message] {
	return retry{onProduceErrors: onProduceErrors, producer: producer}
}
