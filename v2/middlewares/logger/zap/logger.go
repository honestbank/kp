package zap

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/honestbank/kp/v2/middlewares"
	"go.uber.org/zap"
)

type logger struct {
	logger *zap.Logger
}

func (r logger) Process(ctx context.Context, item *kafka.Message, next func(ctx context.Context, item *kafka.Message) error) error {
	scopedLogger := r.logger.
		With(zap.String("key", string(item.Key))).
		With(zap.String("offset", item.TopicPartition.Offset.String())).
		With(zap.Int32("partition", item.TopicPartition.Partition)).
		With(zap.String("topic", getTopicName(item)))

	scopedLogger.Info("Starting to process the message")
	err := next(context.WithValue(ctx, &ctxKey{}, scopedLogger), item)
	if err != nil {
		scopedLogger.Error("Failed while processing the message", zap.Error(err))
	} else {
		scopedLogger.Info("Successfully processed the message")
	}

	return err
}

func NewLoggerMiddleware(zapLogger *zap.Logger) middlewares.KPMiddleware[*kafka.Message] {
	return &logger{logger: zapLogger}
}

func getTopicName(item *kafka.Message) string {
	if item != nil && item.TopicPartition.Topic != nil {
		return *item.TopicPartition.Topic
	}

	return ""
}
