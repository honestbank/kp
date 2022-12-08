package consumer

import "github.com/confluentinc/confluent-kafka-go/kafka"

type Consumer interface {
	GetMessage() *kafka.Message
	Commit(message *kafka.Message) error
}
