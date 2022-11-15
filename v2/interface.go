package v2

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/honestbank/kp/v2/internal/middleware"
)

type KafkaProcessor[MessageType any] interface {
	WithRetry(retryTopic string, retryCount int) (KafkaProcessor[MessageType], error)
	AddMiddleware(middleware middleware.Middleware[*kafka.Message, error]) KafkaProcessor[MessageType]
	WithRetryOrPanic(retryTopic string, retryCount int) KafkaProcessor[MessageType]
	WithDeadletter(deadLetterTopic string) (KafkaProcessor[MessageType], error)
	WithDeadletterOrPanic(deadletterTopic string) KafkaProcessor[MessageType]
	OnKafkaErrors(cb func(err error))
	Stop()
	Run(processor Processor[MessageType]) error
}
