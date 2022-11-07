package producer

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/honestbank/kp/v2/internal/config"
	"github.com/honestbank/kp/v2/internal/schemaregistry"
	"github.com/honestbank/kp/v2/internal/serialization"
)

type producer[BodyType any] struct {
	schemaID int
	k        *kafka.Producer
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

	return p.ProduceRaw(msg)
}

func (p producer[BodyType]) Flush() error {
	// I just felt like 3000 because 3 second should be enough for messages to be flushed
	// we'll need to optimize this as we go
	p.k.Flush(3000)

	return nil
}

func (p producer[BodyType]) ProduceRaw(message *kafka.Message) error {
	// todo: maybe rename this method so that it tells people to not use it unless they know what they're doing.
	message.TopicPartition.Topic = &p.topic
	message.TopicPartition.Partition = kafka.PartitionAny
	return p.k.Produce(message, nil)
}

func New[MessageType any](topic string) (Producer[MessageType], error) {
	cfg, err := config.LoadConfig[config.KafkaConfig]()
	if err != nil {
		return nil, err
	}
	schemaID, err := schemaregistry.Publish[MessageType](topic)
	if err != nil {
		return nil, err
	}

	k, err := kafka.NewProducer(config.GetKafkaConfig(*cfg))

	return producer[MessageType]{
		k:        k,
		schemaID: *schemaID,
		topic:    topic,
	}, nil
}
