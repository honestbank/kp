package kp_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp"
	"github.com/honestbank/kp/mocks"
)

func TestNewConsumer(t *testing.T) {
	t.Run("test consumer", func(t *testing.T) {
		a := assert.New(t)

		ctrl := gomock.NewController(t)
		producer := mocks.NewMockKPProducer(ctrl)

		consumer := kp.NewConsumer("test", "retry-test", "dead-test", 10, func(key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error {
			return nil
		}, producer, 0)

		a.NotNil(consumer)
	})

	t.Run("GetReady", func(t *testing.T) {
		a := assert.New(t)

		ctrl := gomock.NewController(t)
		producer := mocks.NewMockKPProducer(ctrl)

		consumer := kp.NewConsumer("test", "retry-test", "dead-test", 10, func(key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error {
			return nil
		}, producer, 0)

		a.NotNil(consumer)
		a.NotNil(consumer.GetReady())
	})

	t.Run("SetReady", func(t *testing.T) {
		a := assert.New(t)

		ctrl := gomock.NewController(t)
		producer := mocks.NewMockKPProducer(ctrl)

		consumer := kp.NewConsumer("test", "retry-test", "dead-test", 10, func(key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error {
			return nil
		}, producer, 0)

		a.NotNil(consumer)
		a.NotNil(consumer.GetReady())
		consumer.SetReady(make(chan bool))
		a.NotNil(consumer.GetReady())

		consumer.SetReady(nil)
		a.Nil(consumer.GetReady())
	})

	t.Run("Setup", func(t *testing.T) {
		a := assert.New(t)

		ctrl := gomock.NewController(t)
		producer := mocks.NewMockKPProducer(ctrl)

		consumer := kp.NewConsumer("test", "retry-test", "dead-test", 10, func(key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error {
			return nil
		}, producer, 0)

		a.NotNil(consumer)
		a.NotNil(consumer.GetReady())
		a.NotNil(consumer.Setup)
		a.NoError(consumer.Setup(nil))
	})

	t.Run("ProcessWithBackoff", func(t *testing.T) {
		a := assert.New(t)

		ctrl := gomock.NewController(t)
		producer := mocks.NewMockKPProducer(ctrl)

		consumer := kp.NewConsumer("test", "retry-test", "dead-test", 10, func(key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error {
			return nil
		}, producer, time.Second*1)

		message := sarama.ConsumerMessage{
			Headers:        nil,
			Timestamp:      time.Time{},
			BlockTimestamp: time.Time{},
			Key:            nil,
			Value:          []byte(sarama.StringEncoder("test")),
			Topic:          "",
			Partition:      0,
			Offset:         0,
		}
		err := consumer.ProcessWithBackoff(&message)
		a.NoError(err)
	})
	t.Run("ProcessWithBackoff - fail message", func(t *testing.T) {
		a := assert.New(t)

		ctrl := gomock.NewController(t)
		producer := mocks.NewMockKPProducer(ctrl)
		producer.EXPECT().ProduceMessage("retry-test", gomock.Any(), gomock.Any()).Return(nil)

		consumer := kp.NewConsumer("test", "retry-test", "dead-test", 10, func(key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error {
			if message == "fail" {
				return errors.New("fail this message")
			}

			return nil
		}, producer, time.Second*1)

		message := sarama.ConsumerMessage{
			Headers:        nil,
			Timestamp:      time.Time{},
			BlockTimestamp: time.Time{},
			Key:            nil,
			Value:          []byte(sarama.StringEncoder("fail")),
			Topic:          "",
			Partition:      0,
			Offset:         0,
		}

		err := consumer.ProcessWithBackoff(&message)
		a.NoError(err)
	})

	t.Run("ProcessWithBackoff - retry added", func(t *testing.T) {
		a := assert.New(t)

		ctrl := gomock.NewController(t)
		producer := mocks.NewMockKPProducer(ctrl)
		producer.EXPECT().ProduceMessage("retry-test", gomock.Any(), gomock.Any()).Return(nil)

		consumer := kp.NewConsumer("test", "retry-test", "dead-test", 10, func(key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error {
			if message == "fail" {
				return errors.New("fail this message")
			}

			return nil
		}, producer, time.Second*1)

		message := sarama.ConsumerMessage{
			Headers:        nil,
			Timestamp:      time.Time{},
			BlockTimestamp: time.Time{},
			Key:            nil,
			Value:          []byte(sarama.StringEncoder("fail|1")),
			Topic:          "",
			Partition:      0,
			Offset:         0,
		}

		err := consumer.ProcessWithBackoff(&message)
		a.NoError(err)
	})

	t.Run("ProcessWithBackoff - retry exceeded, should not retry", func(t *testing.T) {
		a := assert.New(t)

		ctrl := gomock.NewController(t)
		producer := mocks.NewMockKPProducer(ctrl)
		producer.EXPECT().ProduceMessage("dead-test", gomock.Any(), gomock.Any()).Return(nil)

		consumer := kp.NewConsumer("test", "retry-test", "dead-test", 10, func(key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error {
			if message == "fail" {
				return errors.New("fail this message")
			}

			return nil
		}, producer, time.Second*1)

		message := sarama.ConsumerMessage{
			Headers:        nil,
			Timestamp:      time.Time{},
			BlockTimestamp: time.Time{},
			Key:            nil,
			Value:          []byte(sarama.StringEncoder("fail|10")),
			Topic:          "",
			Partition:      0,
			Offset:         0,
		}

		err := consumer.ProcessWithBackoff(&message)
		a.NoError(err)
	})
}
