package main

import (
	"errors"
	"log"

	"github.com/honestbank/kp"
	"github.com/honestbank/kp/examples/simple/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	// eventing.Setup(*cfg)
	processor := kp.NewKafkaProcessor("test", "dead-test", 10, &cfg.KafkaConfig)
	processor.Process(func(message string) error {
		if message == "fail" {
			return errors.New("failed")
		}
		log.Println("message content:" + message)
		return nil
	})

	processor.Start()

}
