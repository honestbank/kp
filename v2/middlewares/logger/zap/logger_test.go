package zap_test

import (
	"context"
	"errors"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/honestbank/kp/v2/middlewares/logger/zap"
	"github.com/stretchr/testify/assert"
)

func TestLogger_Process(t *testing.T) {
	t.Run("calls next", func(t *testing.T) {
		topicName := "my-topic"
		mockKafkaMessage := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topicName, Offset: kafka.Offset(20), Partition: 2},
			Key:            []byte("my-key"),
		}
		called := false
		zap.NewLoggerMiddleware(zap.LoggerFromContext(context.Background())).Process(context.Background(), mockKafkaMessage, func(ctx context.Context, item *kafka.Message) error {
			called = true

			return nil
		})
		assert.True(t, called)
	})
	t.Run("returns when error is returned", func(t *testing.T) {
		topicName := "my-topic"
		mockKafkaMessage := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topicName, Offset: kafka.Offset(20), Partition: 2},
			Key:            []byte("my-key"),
		}
		called := false
		zap.NewLoggerMiddleware(zap.LoggerFromContext(context.Background())).Process(context.Background(), mockKafkaMessage, func(ctx context.Context, item *kafka.Message) error {
			called = true

			return errors.New("custom error")
		})
		assert.True(t, called)
	})
}
