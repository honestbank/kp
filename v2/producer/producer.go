package producer

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/honestbank/kp/v2/config"
	"github.com/honestbank/kp/v2/internal/schemaregistry"
	"github.com/honestbank/kp/v2/internal/serialization"
	"github.com/honestbank/kp/v2/internal/tracing"
)

type contextKey string

// MessageKey can be used for passing keys before [this](https://github.com/honestbank/kp/issues/81) gets implemented.
const MessageKey contextKey = "messageKey"

type producer[BodyType any] struct {
	schemaID int
	k        UntypedProducer
	topic    string
}

func (p producer[BodyType]) Produce(ctx context.Context, message BodyType) error {
	value, err := serialization.Encode(message, p.schemaID)
	if err != nil {
		return err
	}

	partition := kafka.TopicPartition{
		Topic:     &p.topic,
		Partition: kafka.PartitionAny,
	}
	msg := &kafka.Message{
		TopicPartition: partition,
		Value:          value,
	}
	if val := ctx.Value(MessageKey); val != nil {
		msg.Key = val.([]byte)
	}
	tracing.InjectTraceHeaders(ctx, msg)

	return p.k.ProduceRaw(msg)
}

func (p producer[BodyType]) Flush() error {
	p.k.Flush()

	return nil
}

func (p producer[BodyType]) Events() <-chan kafka.Event {
	return p.k.Events()
}

func (p producer[BodyType]) ProduceRaw(message *kafka.Message) error {
	return p.k.ProduceRaw(message)
}

func New[MessageType any](topic string, cfg config.KPConfig) (Producer[MessageType], error) {
	schemaID, err := schemaregistry.Publish[MessageType](topic, cfg.SchemaRegistryConfig)
	if err != nil {
		return nil, err
	}

	k, err := NewUntyped(topic, cfg.KafkaConfig)
	if err != nil {
		return nil, err
	}

	return producer[MessageType]{
		k:        k,
		schemaID: *schemaID,
		topic:    topic,
	}, nil
}
