package producer

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/honestbank/kp/v2/config"
)

type untypedProducer struct {
	producer *kafka.Producer
	topic    string
}

func (u untypedProducer) ProduceRaw(message *kafka.Message) error {
	message.TopicPartition.Topic = &u.topic
	message.TopicPartition.Partition = kafka.PartitionAny

	return u.producer.Produce(message, nil)
}

func (u untypedProducer) Flush() error {
	// I just felt like 3000 because 3 second should be enough for messages to be flushed
	// we'll need to optimize this as we go
	u.producer.Flush(3_000)

	return nil
}

func NewUntyped(topic string, cfg config.Kafka) (UntypedProducer, error) {
	p, err := kafka.NewProducer(config.GetKafkaConfig(cfg))
	if err != nil {
		return nil, err
	}
	return untypedProducer{
		producer: p,
		topic:    topic,
	}, nil
}
