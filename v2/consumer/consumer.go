package consumer

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/honestbank/kp/v2/config"
)

type consumer struct {
	kakfaConsumer *kafka.Consumer
}

func (c consumer) GetMessage() *kafka.Message {
	ev := c.kakfaConsumer.Poll(100) // kafka's example uses 100ms and I'm going with it for now
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
	_, err := c.kakfaConsumer.CommitMessage(message)

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

	return consumer{kakfaConsumer: k}, nil
}

func NewFromAssignments(topicPartitions []kafka.TopicPartition, cfg config.Kafka) (Consumer, error) {
	kafkaConfig := config.GetKafkaConsumerConfig(cfg)
	_ = kafkaConfig.SetKey("enable.auto.commit", false)
	k, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		return nil, err
	}

	err = k.Assign(topicPartitions)
	if err != nil {
		return nil, err
	}
	return consumer{kakfaConsumer: k}, nil
}

func (c consumer) GetAssignments() ([]kafka.TopicPartition, error) {
	return c.kakfaConsumer.Assignment()
}
