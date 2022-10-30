package producer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/honestbank/kp/v2/internal/config"
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

	return p.ProduceRaw(msg)
}

func (p producer[BodyType, KeyType]) Flush() error {
	// I just felt like 3000 because 3 second should be enough for messages to be flushed
	// we'll need to optimize this as we go
	p.k.Flush(3000)

	return nil
}

func (p producer[BodyType, KeyType]) ProduceRaw(message *kafka.Message) error {
	// todo: maybe rename this method so that it tells people to not use it unless they know what they're doing.
	message.TopicPartition.Topic = &p.topic
	message.TopicPartition.Partition = kafka.PartitionAny
	return p.k.Produce(message, nil)
}

func New[MessageType any, KeyType KeyTypes](topic string) (Producer[MessageType, KeyType], error) {
	cfg, err := config.LoadConfig[config.KafkaConfig]()
	if err != nil {
		return nil, err
	}
	schemaID, err := schemaregistry.Publish[MessageType](topic)
	if err != nil {
		return nil, err
	}

	k, err := kafka.NewProducer(config.GetKafkaConfig(*cfg))

	return producer[MessageType, KeyType]{
		k:        k,
		schemaID: *schemaID,
		topic:    topic,
	}, nil
}
