package kp

import "github.com/Shopify/sarama"

type KafkaProcessor interface {
	Process(processor func(key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error)
	Start()
	Stop()
	OnFailure(failure func(key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error)
}
