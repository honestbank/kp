package consumer

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/honestbank/kp/v2/config"
)

type consumer struct {
	kafkaConsumer *kafka.Consumer
}

// TopicPartitionOffset can be used to initialize a Consumer that can start consuming from the given offset on a given partition in the given topic
type TopicPartitionOffset struct {
	Topic     *string
	Partition int32
	Offset    int64
}

func (c consumer) GetMessage() *kafka.Message {
	ev := c.kafkaConsumer.Poll(100) // kafka's example uses 100ms and I'm going with it for now
	if ev == nil {
		return nil
	}

	return getMessageOrNil(ev)
}

func getMessageOrNil(event kafka.Event) *kafka.Message {
	switch e := event.(type) {
	case *kafka.Message:
		return e
	default:
		// Errors should generally be considered informational, the client will try to automatically recover.
		return nil
	}
}

func (c consumer) Commit(message *kafka.Message) error {
	_, err := c.kafkaConsumer.CommitMessage(message)

	return err
}

func New(topics []string, cfg config.Kafka) (Consumer, error) {
	kafkaConfig := config.GetKafkaConsumerConfig(cfg)
	_ = kafkaConfig.SetKey("enable.auto.commit", false)
	k, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		return nil, err
	}
	err = k.SubscribeTopics(topics, nil)
	if err != nil {
		return nil, err
	}

	return consumer{kafkaConsumer: k}, nil
}

func NewFromAssignments(topicPartitionOffsets []TopicPartitionOffset, cfg config.Kafka) (Consumer, error) {
	kafkaConfig := config.GetKafkaConsumerConfig(cfg)
	_ = kafkaConfig.SetKey("enable.auto.commit", false)
	k, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		return nil, err
	}

	var assignments []kafka.TopicPartition
	for _, topicPartitionOffset := range topicPartitionOffsets {
		assignments = append(assignments, kafka.TopicPartition{
			Topic:     topicPartitionOffset.Topic,
			Partition: topicPartitionOffset.Partition,
			Offset:    kafka.Offset(topicPartitionOffset.Offset),
		})
	}

	err = k.Assign(assignments)
	if err != nil {
		return nil, err
	}
	return consumer{kafkaConsumer: k}, nil
}
