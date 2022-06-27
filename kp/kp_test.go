package kp_test

import (
	"errors"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/kp"
)

var kafkaConfig = kp.KafkaConfig{
	KafkaBootstrap: "localhost:9092",
}

var producer = kp.NewProducer(kafkaConfig)

func TestNewKafkaProcessor(t *testing.T) {
	t.Run("test new kafka processor", func(t *testing.T) {
		a := assert.New(t)

		processor := kp.NewKafkaProcessor("test", "dead-test", 10, kafkaConfig)
		a.NotNil(processor)
	})

	t.Run("test process", func(t *testing.T) {
		a := assert.New(t)
		data := make([]string, 0)

		processor := kp.NewKafkaProcessor("test", "dead-test", 10, kafkaConfig)
		processor.Process(func(key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error {
			data = append(data, message)
			if message == "fail" {
				return errors.New("failed")
			}
			log.Println("message content:" + message)

			return nil
		})

		var wg sync.WaitGroup
		wg.Add(1)
		quit := make(chan bool)

		go func() {
			for {
				foo, ok := <-quit
				if !ok {
					processor.Stop()
					wg.Done()

					return
				}
				log.Println("foo:", foo)
				processor.Start()
				// Do other stuff
			}
		}()
		quit <- true
		time.Sleep(time.Second * 5)

		err := producer.ProduceMessage("test", "1", "test")
		a.NoError(err)
		err = producer.ProduceMessage("test", "2", "test")
		a.NoError(err)
		err = producer.ProduceMessage("test", "3", "test")
		a.NoError(err)
		time.Sleep(time.Second * 5)
		close(quit)
		a.Equal(3, len(data))
	})
	t.Run("test process fail", func(t *testing.T) {
		a := assert.New(t)
		data := make([]string, 0)

		processor := kp.NewKafkaProcessor("test-fail", "dead-test-fail", 10, kafkaConfig)
		processor.Process(func(key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error {
			data = append(data, message)
			if message == "fail" {
				return errors.New("failed")
			}
			log.Println("message content:" + message)

			return nil
		})

		var wg sync.WaitGroup
		wg.Add(1)
		quit := make(chan bool)

		go func() {
			for {
				foo, ok := <-quit
				if !ok {
					processor.Stop()
					wg.Done()

					return
				}
				log.Println("foo:", foo)
				processor.Start()
				// Do other stuff
			}
		}()
		quit <- true
		time.Sleep(time.Second * 5)

		_ = producer.ProduceMessage("test-fail", "1", "fail")
		time.Sleep(time.Second * 5)
		close(quit)
		a.Equal(10, len(data))
	})
}
