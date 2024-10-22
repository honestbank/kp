package producer

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/honestbank/kp/v2/internal/middleware"
)

type Producer[BodyType any] interface {
	Flush() error
	Produce(context context.Context, message BodyType) error
	SetMiddlewares(middlewares []middleware.Middleware[*kafka.Message, error])
}

type UntypedProducer interface {
	Produce(ctx context.Context, message *kafka.Message) error
	SetMiddlewares(middlewares []middleware.Middleware[*kafka.Message, error])
	Flush() error
}
