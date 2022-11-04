package v2

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/honestbank/kp/v2/internal/middleware"
)

type KafkaProcessor[MessageType any] interface {
	WithRetry(retryTopic string, retryCount int) (KafkaProcessor[MessageType], error)
	AddMiddleware(middleware middleware.Middleware[*kafka.Message, error]) KafkaProcessor[MessageType]
	WithRetryOrPanic(retryTopic string, retryCount int) KafkaProcessor[MessageType]
	WithDeadletter(deadLetterTopic string) (KafkaProcessor[MessageType], error)
	WithDeadletterOrPanic(deadletterTopic string) KafkaProcessor[MessageType]
	Stop()
	Run(processor func(ctx context.Context, message MessageType) error) error
}
