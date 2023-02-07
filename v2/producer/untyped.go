package producer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"

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
	u.producer.Flush(3_000)

	return nil
}

func NewUntyped(topic string, cfg config.KPConfig) (UntypedProducer, error) {
	p, err := kafka.NewProducer(config.GetKafkaConfig(cfg.KafkaConfig))
	if err != nil {
		return nil, err
	}
	return untypedProducer{
		producer: p,
		topic:    topic,
	}, nil
}
