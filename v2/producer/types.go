package producer

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer[BodyType any] interface {
	Flush() error
	Produce(context context.Context, message BodyType) error
	ProduceRaw(message *kafka.Message) error
}
