package producer

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/honestbank/kp/v2/config"
	"github.com/honestbank/kp/v2/internal/schemaregistry"
	"github.com/honestbank/kp/v2/internal/serialization"
	"github.com/honestbank/kp/v2/internal/tracing"
)

type producer[BodyType any] struct {
	schemaID int
	k        UntypedProducer
	topic    string
}

func (p producer[BodyType]) Produce(ctx context.Context, message BodyType) error {
	return doProduce(ctx, p.k, nil, message, &p.schemaID, &p.topic)
}

func (p producer[BodyType]) ProduceWithKey(ctx context.Context, key []byte, message BodyType) error {
	return doProduce(ctx, p.k, key, message, &p.schemaID, &p.topic)
}

func doProduce[BodyType any](ctx context.Context, p UntypedProducer, key []byte, message BodyType, schemaID *int, topic *string) error {
	value, err := serialization.Encode(message, *schemaID)
	if err != nil {
		return err
	}

	partition := kafka.TopicPartition{
		Topic:     topic,
		Partition: kafka.PartitionAny,
	}

	msg := &kafka.Message{
		Key:            key,
		TopicPartition: partition,
		Value:          value,
	}

	tracing.InjectTraceHeaders(ctx, msg)

	return p.ProduceRaw(msg)
}

func (p producer[BodyType]) Flush() error {
	p.k.Flush()

	return nil
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
