package producer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/honestbank/kp/v2/internal/schemaregistry"
	"github.com/honestbank/kp/v2/internal/serialization"
)

type producer[BodyType any, KeyType KeyTypes] struct {
	schemaID int
	k        *kafka.Producer
	topic    string
}

func (p producer[BodyType, KeyType]) Produce(message KafkaMessage[BodyType, KeyType]) error {
	value, err := serialization.Encode(message.Body, p.schemaID)
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

	return p.k.Produce(msg, nil)
}

func (p producer[BodyType, KeyType]) Flush() error {
	p.k.Flush(3000)

	return nil
}

func New[MessageType any, KeyType KeyTypes](topic string) (Producer[MessageType, KeyType], error) {
	schemaID, err := schemaregistry.Publish[MessageType](topic)
	if err != nil {
		return nil, err
	}

	k, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})

	return producer[MessageType, KeyType]{
		k:        k,
		schemaID: *schemaID,
		topic:    topic,
	}, nil
}
