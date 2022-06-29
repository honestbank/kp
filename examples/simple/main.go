package main

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/Shopify/sarama"

	"github.com/honestbank/kp"
	"github.com/honestbank/kp/examples/simple/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	processor := kp.NewKafkaProcessor("test", "retry-test", "dead-test", 10, "simple-service", kp.KafkaConfig{KafkaBootstrapServers: strings.Split(cfg.KafkaConfig.KafkaBootstrapServers, ",")}, 0)
	processor.Process(func(ctx context.Context, key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error {
		if message == "fail" {
			return errors.New("failed")
		}
		log.Println("message content:" + message)
		return nil
	})

	processor.Start(ctx)

}
