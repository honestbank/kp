package v2_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/assert"

	v2 "github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/config"
	consumer2 "github.com/honestbank/kp/v2/consumer"
	"github.com/honestbank/kp/v2/middlewares/consumer"
	"github.com/honestbank/kp/v2/middlewares/deadletter"
	"github.com/honestbank/kp/v2/middlewares/retry"
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
	kafkaCfg := config.Kafka{BootstrapServers: "localhost", ConsumerGroupName: "integration-tests"}
	schemaRegistryConfig := config.SchemaRegistry{Endpoint: "http://localhost:8082"}
	c, err := kafka.NewAdminClient(config.GetKafkaConfig(config.Kafka{BootstrapServers: "localhost"}))
	assert.NoError(t, err)
	_, err = c.CreateTopics(context.Background(), []kafka.TopicSpecification{{Topic: "kp-topic", ReplicationFactor: 1, NumPartitions: 1}, {Topic: "kp-topic-retry", ReplicationFactor: 1, NumPartitions: 1}, {Topic: "kp-topic-dlt", ReplicationFactor: 1, NumPartitions: 1}})
	assert.NoError(t, err)

	p, err := producer.New[MyType]("kp-topic", config.KPConfig{KafkaConfig: kafkaCfg, SchemaRegistryConfig: schemaRegistryConfig})
	p.Produce(context.Background(), MyType{Username: "username1", Count: 1})
	p.Flush()
	assert.NoError(t, err)
	kp := v2.New[kafka.Message](nil)
	messageProcessCount := 0
	const retryCount = 10
	retryTopicProducer, err := producer.New[UserLoggedInEvent]("kp-topic-retry", config.KPConfig{KafkaConfig: kafkaCfg, SchemaRegistryConfig: schemaRegistryConfig})
	if err != nil {
		panic(err)
	}
	dltProducer, err := producer.New[UserLoggedInEvent]("kp-topic-dlt", config.KPConfig{KafkaConfig: kafkaCfg, SchemaRegistryConfig: schemaRegistryConfig})
	if err != nil {
		panic(err)
	}
	kafkaConsumer, err := consumer2.New([]string{"kp-topic", "kp-topic-retry"}, kafkaCfg.WithDefaults())
	if err != nil {
		panic(err)
	}

	go func() {
		time.Sleep(time.Second * (retryCount / 2))
		kp.Stop()
	}()
	err = kp.AddMiddleware(consumer.NewConsumerMiddleware(kafkaConsumer)).
		AddMiddleware(MyMw{}).
		AddMiddleware(retry.NewRetryMiddleware(retryTopicProducer)).
		AddMiddleware(deadletter.NewDeadletterMiddleware(dltProducer, retryCount)).
		Run(func(ctx context.Context, message *kafka.Message) error {
			fmt.Printf("%v\n", message)
			messageProcessCount++

			return errors.New("some error")
		})
	assert.Equal(t, retryCount+1, messageProcessCount)
	assert.NoError(t, err)
}
