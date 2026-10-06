package producer

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer[BodyType any] interface {
	Close()
	Events() <-chan kafka.Event
	Flush() error
	Produce(context context.Context, message BodyType) error
	ProduceRaw(message *kafka.Message) error
}

type UntypedProducer interface {
	Close()
	Events() <-chan kafka.Event
	Flush() error
	ProduceRaw(message *kafka.Message) error
}
