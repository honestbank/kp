package v2_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	v2 "github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/config"
	consumer2 "github.com/honestbank/kp/v2/consumer"
	"github.com/honestbank/kp/v2/internal/serialization"
	"github.com/honestbank/kp/v2/middlewares/consumer"
	"github.com/honestbank/kp/v2/middlewares/deadletter"
	"github.com/honestbank/kp/v2/middlewares/retry"
	"github.com/honestbank/kp/v2/middlewares/retry_count"
	"github.com/honestbank/kp/v2/producer"
)

type UserLoggedInEvent struct {
	UserID    string
	Timestamp string
}

func ExampleNew() {
	setup()
	kpConfig := config.KPConfig{
		KafkaConfig: config.Kafka{
			BootstrapServers:  "localhost",
			ConsumerGroupName: "example-tests",
		},
		SchemaRegistryConfig: config.SchemaRegistry{
			Endpoint: "http://localhost:8082",
		},
	}

	retryTopicProducer, err := producer.New[UserLoggedInEvent]("user-logged-in-rewards-processor-retry", kpConfig)
	if err != nil {
		panic(err)
	}
	dltProducer, err := producer.New[UserLoggedInEvent]("user-logged-in-rewards-processor-dlt", kpConfig)
	if err != nil {
		panic(err)
	}
	processor := v2.New[kafka.Message](nil)
	kafkaConsumer, err := consumer2.New([]string{"user-logged-in", "user-logged-in-rewards-processor-retry"}, kpConfig.KafkaConfig.WithDefaults())
	if err != nil {
		panic(err)
	}
	go func() {
		time.Sleep(time.Second * 10)
		processor.Stop()
	}()
	err = processor.
		AddMiddleware(consumer.NewConsumerMiddleware(kafkaConsumer)).
		AddMiddleware(retry_count.NewRetryCountMiddleware()).
		AddMiddleware(retry.NewRetryMiddleware(retryTopicProducer)).
		AddMiddleware(deadletter.NewDeadletterMiddleware(dltProducer, 3)).
		Run(func(ctx context.Context, ev *kafka.Message) error {
			val, _ := serialization.Decode[UserLoggedInEvent](ev.Value)
			fmt.Printf("%s-%d|", val.UserID, retry_count.FromContext(ctx))
			return errors.New("some error")
		})
	if err != nil {
		panic(err)
	}
	// Output: 1-0|1-1|1-2|1-3|
}

func setup() {
	cfg := config.KPConfig{KafkaConfig: config.Kafka{BootstrapServers: "localhost"}, SchemaRegistryConfig: config.SchemaRegistry{Endpoint: "http://localhost:8082"}}
	c, err := kafka.NewAdminClient(config.GetKafkaConfig(cfg.KafkaConfig))
	if err != nil {
		panic(err)
	}
	_, err = c.CreateTopics(context.Background(), []kafka.TopicSpecification{{Topic: "user-logged-in", ReplicationFactor: 1, NumPartitions: 1}, {Topic: "user-logged-in-rewards-processor-retry", ReplicationFactor: 1, NumPartitions: 1}, {Topic: "user-logged-in-rewards-processor-dlt", ReplicationFactor: 1, NumPartitions: 1}})
	if err != nil {
		panic(err)
	}
	p, err := producer.New[UserLoggedInEvent]("user-logged-in", cfg)
	if err != nil {
		panic(err)
	}
	now := time.Now().Format(time.RFC3339Nano)
	event := UserLoggedInEvent{UserID: "1", Timestamp: now}
	err = p.Produce(context.Background(), event)
	if err != nil {
		panic(err)
	}
	err = p.Flush()
	if err != nil {
		panic(err)
	}
}
