package kp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/kp"
)

var kafkaConfig = kp.KafkaConfig{
	KafkaBootstrap: "localhost:9092",
}

func TestNewKafkaProcessor(t *testing.T) {
	t.Run("test new kafka processor", func(t *testing.T) {
		a := assert.New(t)

		processor := kp.NewKafkaProcessor("test", "dead-test", 10, kafkaConfig)
		a.NotNil(processor)
	})
}
