//go:build lab_integration_test

package producer_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/producer"
)

type BenchmarkMessage struct {
	Body  string `json:"body" avro:"body"`
	Count int    `json:"count" avro:"count"`
}

func TestNewCanConnectToLabEnvironment(t *testing.T) {
	_ = os.Setenv("KP_SASL_MECHANISM", "PLAIN")
	_ = os.Setenv("KP_SECURITY_PROTOCOL", "SASL_SSL")
	defer func() {
		_ = os.Unsetenv("KP_SASL_MECHANISM")
		_ = os.Unsetenv("KP_SECURITY_PROTOCOL")
	}()
	t.Run("produce through kp", func(t *testing.T) {
		kp, err := producer.New[BenchmarkMessage, int]("topic-kp")
		assert.NoError(t, err)
		defer kp.Flush()

		for i := 0; i < 100; i++ {
			err := kp.Produce(producer.KafkaMessage[BenchmarkMessage, int]{
				Body: BenchmarkMessage{
					Body:  "hello-world",
					Count: i,
				},
			})
			assert.NoError(t, err)
		}
	})
}
