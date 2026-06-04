package consumer

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/honestbank/kp/v2/consumer"
	"github.com/honestbank/kp/v2/middlewares"
)

type consumerMiddleware struct {
	consumer consumer.Consumer
}

func (c consumerMiddleware) Process(ctx context.Context, item *kafka.Message, next func(ctx context.Context, item *kafka.Message) error) (retErr error) {
	if item != nil {
		// I don't think this will ever happen though...
		return next(ctx, item)
	}
	msg := c.consumer.GetMessage()
	if msg == nil {
		return nil
	}
	defer func() {
		if err := c.consumer.Commit(msg); err != nil && retErr == nil {
			retErr = err
		}
	}()

	return next(ctx, msg)
}

func NewConsumerMiddleware(consumer consumer.Consumer) middlewares.KPMiddleware[*kafka.Message] {
	return &consumerMiddleware{
		consumer: consumer,
	}
}
