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

func (c consumerMiddleware) Process(ctx context.Context, item *kafka.Message, next func(ctx context.Context, item *kafka.Message) error) error {
	if item != nil {
		// I don't think this will ever happen though...
		return next(ctx, item)
	}
	msg := c.consumer.GetMessage()
	if msg == nil {
		return nil
	}
	defer c.consumer.Commit(msg)

	return next(ctx, msg)
}

// Close leaves the consumer group and releases the underlying handle. kp.Run
// invokes it automatically once the processing loop returns, via the io.Closer
// interface, so downstream services get graceful group-leave from a version
// bump alone with no manual defer.
func (c consumerMiddleware) Close() error {
	return c.consumer.Close()
}

func NewConsumerMiddleware(consumer consumer.Consumer) middlewares.KPMiddleware[*kafka.Message] {
	return &consumerMiddleware{
		consumer: consumer,
	}
}
