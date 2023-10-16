package producer

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer[BodyType any] interface {
	Flush() error
	Produce(context context.Context, message BodyType) error
	ProduceWithKey(context context.Context, key []byte, message BodyType) error
	ProduceRaw(message *kafka.Message) error
}

type UntypedProducer interface {
	ProduceRaw(message *kafka.Message) error
	Flush() error
}
