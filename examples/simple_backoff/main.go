package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Shopify/sarama"

	"github.com/honestbank/kp"
	"github.com/honestbank/kp/examples/simple/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	processor := kp.NewKafkaProcessor("test", "retry-test", "dead-test", 10, "simple-service", kp.KafkaConfig{KafkaBootstrapServers: strings.Split(cfg.KafkaConfig.KafkaBootstrapServers, ",")}, time.Second*2)
	processor.Process(func(ctx context.Context, key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error {
		time.Sleep(time.Second * 60)
		if message == "fail" {
			return errors.New("failed")
		}
		return nil
	})

	processor.Start(context.Background())

}
