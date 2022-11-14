package v2_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	v2 "github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/config"
	"github.com/honestbank/kp/v2/producer"
)

type UserLoggedInEvent struct {
	UserID    string
	Timestamp string
}

func ExampleNew() {
	setup()

	processor := v2.New[UserLoggedInEvent]("user-logged-in", config.KPConfig{KafkaConfig: config.Kafka{BootstrapServers: "localhost"}, SchemaRegistryConfig: config.SchemaRegistry{Endpoint: "http://localhost:8081"}})
	go func() {
		time.Sleep(time.Second * 10)
		processor.Stop()
	}()
	processor.WithRetryOrPanic("user-logged-in-rewards-processor-retry", 3).
		WithDeadletterOrPanic("user-logged-in-rewards-processor-dlt").
		Run(func(ctx context.Context, ev UserLoggedInEvent) error {
			fmt.Printf("%s|", ev.UserID)
			return errors.New("some error")
		})
	// Output: 1|1|1|1|
}

func setup() {
	cfg := config.KPConfig{KafkaConfig: config.Kafka{BootstrapServers: "localhost"}, SchemaRegistryConfig: config.SchemaRegistry{Endpoint: "http://localhost:8081"}}
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
	p.Produce(context.Background(), event)
	p.Flush()
}
