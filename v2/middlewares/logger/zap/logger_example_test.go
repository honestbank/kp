package zap_test

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/honestbank/kp/v2/middlewares/logger/zap"
	uberzap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ExampleNewLoggerMiddleware() {
	mw := zap.NewLoggerMiddleware(getInitialLogger())

	topicName := "my-topic"
	mockKafkaMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicName, Offset: kafka.Offset(20), Partition: 2},
		Key:            []byte("my-key"),
	}
	_ = mw.Process(context.Background(), mockKafkaMessage, func(ctx context.Context, item *kafka.Message) error {
		logger := zap.LoggerFromContext(ctx)
		logger.Info("processing message")
		logger.Error("returning an error")

		return errors.New("custom_error")
	})

	// output: {"level":"INFO","ts":"<timestamp>","msg":"Starting to process the message","key":"my-key","offset":"20","partition":2,"topic":"my-topic"}
	//{"level":"INFO","ts":"<timestamp>","msg":"processing message","key":"my-key","offset":"20","partition":2,"topic":"my-topic"}
	//{"level":"ERROR","ts":"<timestamp>","msg":"returning an error","key":"my-key","offset":"20","partition":2,"topic":"my-topic"}
	//{"level":"ERROR","ts":"<timestamp>","msg":"Failed while processing the message","key":"my-key","offset":"20","partition":2,"topic":"my-topic","error":"custom_error"}
}

func getInitialLogger() *uberzap.Logger {
	encoderCfg := uberzap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("<timestamp>")
	encoderCfg.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(strings.ToUpper(level.String()))
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		uberzap.NewAtomicLevelAt(zapcore.DebugLevel),
	)

	return uberzap.New(core)
}
