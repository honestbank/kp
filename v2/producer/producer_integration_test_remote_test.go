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
	_ = os.Setenv("KP_USERNAME", "OBHGFSK7ZAVNEDFV")
	_ = os.Setenv("KP_PASSWORD", "P5yMMU05T9MNiJ9FHUINoEpm7PeW0huuyuF0H01FEZe4YcDBkdkYRI0nAi/7/bFv")
	_ = os.Setenv("KP_KAFKA_BOOTSTRAP_SERVERS", "pkc-ew3qg.asia-southeast2.gcp.confluent.cloud:9092")
	_ = os.Setenv("KP_SASL_MECHANISM", "PLAIN")
	_ = os.Setenv("KP_SECURITY_PROTOCOL", "SASL_SSL")
	_ = os.Setenv("KP_SCHEMA_REGISTRY_ENDPOINT", "https://psrc-3w372.australia-southeast1.gcp.confluent.cloud")
	_ = os.Setenv("KP_SCHEMA_REGISTRY_USERNAME", "SVZXV26UZ7TGPU7D")
	_ = os.Setenv("KP_SCHEMA_REGISTRY_PASSWORD", "IbnoBsE83SbTqZjHznFcLxFsOtJoWzgw7CZunhVnOmtL8wJzhv9IVaVF0bVzLd9/")
	defer func() {
		_ = os.Unsetenv("KP_USERNAME")
		_ = os.Unsetenv("KP_PASSWORD")
		_ = os.Unsetenv("KP_KAFKA_BOOTSTRAP_SERVERS")
		_ = os.Unsetenv("KP_SCHEMA_REGISTRY_ENDPOINT")
		_ = os.Unsetenv("KP_SCHEMA_REGISTRY_USERNAME")
		_ = os.Unsetenv("KP_SCHEMA_REGISTRY_PASSWORD")
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
