package consumer

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type Consumer interface {
	GetMessage() *kafka.Message
	Commit(message *kafka.Message) error
	Stop() error
}
