package producer

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer[BodyType any] interface {
	Flush() error
	Produce(context context.Context, message BodyType) error
	ProduceRaw(message *kafka.Message) error
	GetTopic() string
}

type UntypedProducer interface {
	ProduceRaw(message *kafka.Message) error
	GetTopic() string
	Flush() error
}
