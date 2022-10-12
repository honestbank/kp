package kp

import (
	"context"

	"github.com/Shopify/sarama"
)

type KafkaProcessor interface {
	Process(processor func(ctx context.Context, key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error)
	Start(ctx context.Context) error
	Stop()
	OnFailure(failure func(ctx context.Context, key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error)
}
