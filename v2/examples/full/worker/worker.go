package worker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	v2 "github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/middlewares"
	"github.com/honestbank/kp_demo/types"
)

func Process() {
	ensureTopicsExist()

	tracing := middlewares.Tracing()
	err := v2.New[types.RegistrationEvent]("user-registration-complete", "send-welcome-email-worker").
		WithRetryOrPanic("send-welcome-email-retries", 10).
		WithDeadletterOrPanic("send-welcome-email-failures").
		AddMiddleware(middlewares.RetryCount()).
		AddMiddleware(tracing).
		Run(func(ctx context.Context, message types.RegistrationEvent) error {
			time.Sleep(time.Millisecond * 300)
			fmt.Printf("retry count: %d", middlewares.RetryCountFromContext(ctx))
			if middlewares.RetryCountFromContext(ctx) <= 5 {
				return errors.New("some error")
			}
			return nil
		})
	if err != nil {
		panic(err)
	}
}

func ensureTopicsExist() {
	cfg := &kafka.ConfigMap{}
	cfg.SetKey("bootstrap.servers", os.Getenv("KP_KAFKA_BOOTSTRAP_SERVERS"))
	// ensure the topics actually exist.
	c, err := kafka.NewAdminClient(cfg)
	if err != nil {
		panic(err)
	}
	_, err = c.CreateTopics(context.Background(), []kafka.TopicSpecification{
		{Topic: "user-registration-complete", ReplicationFactor: 1, NumPartitions: 1},
		{Topic: "send-welcome-email-retries", ReplicationFactor: 1, NumPartitions: 1},
		{Topic: "send-welcome-email-failures", ReplicationFactor: 1, NumPartitions: 1}})
}
