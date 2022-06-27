package main

import (
	"errors"
	"log"

	"github.com/Shopify/sarama"

	"github.com/honestbank/kp"
	"github.com/honestbank/kp/examples/simple/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	// eventing.Setup(*cfg)
	processor := kp.NewKafkaProcessor("test", "dead-test", 10, kp.KafkaConfig{KafkaBootstrapServers: cfg.KafkaConfig.KafkaBootstrapServers})
	processor.Process(func(key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error {
		if message == "fail" {
			return errors.New("failed")
		}
		log.Println("message content:" + message)
		return nil
	})

	processor.Start()

}
