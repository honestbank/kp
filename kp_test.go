package kp_test

import (
	"errors"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"

	kp2 "github.com/honestbank/kp"
)

var kafkaConfig = kp2.KafkaConfig{
	KafkaBootstrapServers: []string{"localhost:9092"},
}

var producer = kp2.NewProducer(kafkaConfig)

type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

// Inc increments the counter for the given key.
func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[key]++
	c.mu.Unlock()
}

// Value returns the current value of the counter for the given key.
func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mu.Unlock()

	return c.v[key]
}

func TestNewKafkaProcessor(t *testing.T) {
	t.Run("test new kafka processor", func(t *testing.T) {
		a := assert.New(t)

		processor := kp2.NewKafkaProcessor("test", "retry-test", "dead-test", 10, "test", kafkaConfig, 0)
		a.NotNil(processor)
	})

	t.Run("test process", func(t *testing.T) {
		a := assert.New(t)
		data := SafeCounter{v: make(map[string]int)}

		processor := kp2.NewKafkaProcessor("test", "retry-test", "dead-test", 10, "test", kafkaConfig, 0)
		processor.Process(func(key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error {
			data.Inc("test")
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
		time.Sleep(time.Second * 10)

		err := producer.ProduceMessage("test", "1", "test")
		a.NoError(err)
		err = producer.ProduceMessage("test", "2", "test")
		a.NoError(err)
		err = producer.ProduceMessage("test", "3", "test")
		a.NoError(err)
		time.Sleep(time.Second * 5)
		close(quit)
		a.Equal(3, data.Value("test"))
	})
	t.Run("test process fail", func(t *testing.T) {
		a := assert.New(t)
		data := SafeCounter{v: make(map[string]int)}

		processor := kp2.NewKafkaProcessor("test-fail", "retry-test-fail", "dead-test-fail", 10, "test-fail", kafkaConfig, 0)
		processor.Process(func(key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error {
			data.Inc("fail")
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
		time.Sleep(time.Second * 10)

		_ = producer.ProduceMessage("test-fail", "1", "fail")
		time.Sleep(time.Second * 5)
		close(quit)
		a.Equal(10, data.Value("fail"))
	})
}
