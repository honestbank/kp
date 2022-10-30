package consumer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/honestbank/kp/v2/internal/config"
)

type consumer struct {
	c *kafka.Consumer
}

func (c consumer) GetMessage() *kafka.Message {
	ev := c.c.Poll(100) // kafka's example uses 100ms and I'm going with it for now
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
	_, err := c.c.CommitMessage(message)

	return err
}

func New(topics []string, consumerGroup string) (Consumer, error) {
	cfg, err := config.LoadConfig[config.KafkaConfig]()
	if err != nil {
		return nil, err
	}
	kafkaConfig := config.GetKafkaConsumerConfig(*cfg)
	_ = kafkaConfig.SetKey("group.id", consumerGroup)
	k, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		return nil, err
	}
	err = k.SubscribeTopics(topics, nil)
	if err != nil {
		return nil, err
	}

	return consumer{c: k}, nil
}
