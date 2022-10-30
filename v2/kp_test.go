package v2_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	v2 "github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/producer"
)

type MyType struct {
	Count    int
	Username string
}

func TestKP(t *testing.T) {
	t.Setenv("KP_SCHEMA_REGISTRY_ENDPOINT", "http://localhost:8081")
	t.Setenv("KP_KAFKA_BOOTSTRAP_SERVERS", "localhost")
	p, err := producer.New[MyType, string]("kp-topic")
	p.Produce(producer.KafkaMessage[MyType, string]{Body: MyType{Username: "username1", Count: 1}})
	p.Flush()
	assert.NoError(t, err)
	kp := v2.New[MyType]("kp-topic")
	messageProcessCount := 0
	const retryCount = 10
	go func() {
		time.Sleep(time.Second * (retryCount / 2))
		kp.Stop()
	}()
	err = kp.WithRetryOrPanic("kp-topic-retry", retryCount).
		WithDeadletterOrPanic("kp-topic-dlt").
		Run(func(message MyType) error {
			time.Sleep(time.Millisecond * 200)
			messageProcessCount++

			return errors.New("some error")
		})
	assert.Equal(t, retryCount+1, messageProcessCount)
	assert.NoError(t, err)
}
