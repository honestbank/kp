package v2_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"

	v2 "github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/internal/config"
	"github.com/honestbank/kp/v2/producer"
)

type MyType struct {
	Count    int
	Username string
}

type MyMw struct {
}

func (m MyMw) Process(ctx context.Context, item *kafka.Message, next func(ctx context.Context, item *kafka.Message) error) error {
	fmt.Println("Before:")
	result := next(ctx, item)
	fmt.Printf("After with return value: %v\n", result)

	return result
}

func TestKP(t *testing.T) {
	t.Setenv("KP_SCHEMA_REGISTRY_ENDPOINT", "http://localhost:8081")
	t.Setenv("KP_KAFKA_BOOTSTRAP_SERVERS", "localhost")
	cfg, err := config.LoadConfig[config.KafkaConfig]()
	assert.NoError(t, err)
	c, err := kafka.NewAdminClient(config.GetKafkaConfig(*cfg))
	assert.NoError(t, err)
	_, err = c.CreateTopics(context.Background(), []kafka.TopicSpecification{{Topic: "kp-topic", ReplicationFactor: 1, NumPartitions: 1}, {Topic: "kp-topic-retry", ReplicationFactor: 1, NumPartitions: 1}, {Topic: "kp-topic-dlt", ReplicationFactor: 1, NumPartitions: 1}})
	assert.NoError(t, err)

	p, err := producer.New[MyType, string]("kp-topic")
	p.Produce(producer.KafkaMessage[MyType, string]{Body: MyType{Username: "username1", Count: 1}})
	p.Flush()
	assert.NoError(t, err)
	kp := v2.New[MyType]("kp-topic", "integration-tests")
	messageProcessCount := 0
	const retryCount = 10
	go func() {
		time.Sleep(time.Second * (retryCount / 2))
		kp.Stop()
	}()
	err = kp.WithRetryOrPanic("kp-topic-retry", retryCount).
		AddMiddleware(MyMw{}).
		WithDeadletterOrPanic("kp-topic-dlt").
		Run(func(ctx context.Context, message MyType) error {
			time.Sleep(time.Millisecond * 200)
			fmt.Printf("%v\n", message)
			messageProcessCount++

			return errors.New("some error")
		})
	assert.Equal(t, retryCount+1, messageProcessCount)
	assert.NoError(t, err)
}
