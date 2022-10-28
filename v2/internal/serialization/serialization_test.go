//go:build integration_test

package serialization_test

import (
	"os"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/avro"
	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/internal/serialization"
)

type BenchmarkMessage struct {
	Body  string `json:"body" avro:"body"`
	Count int    `json:"count" avro:"count"`
}

func TestSerialization(t *testing.T) {
	t.Run("can serialize and deserialize", func(t *testing.T) {
		bytes, err := serialization.Encode(BenchmarkMessage{
			Body:  "my-body",
			Count: 100,
		}, 1)
		assert.NoError(t, err)
		assert.NotNil(t, bytes)
		msg, err := serialization.Decode[BenchmarkMessage](bytes)
		assert.NoError(t, err)
		assert.NotNil(t, msg)
		assert.Equal(t, "my-body", msg.Body)
		assert.Equal(t, 100, msg.Count)
	})
	t.Run("matches with confluent serializer", func(t *testing.T) {
		os.Setenv("KP_SCHEMA_REGISTRY_ENDPOINT", "http://localhost:8081")
		defer os.Unsetenv("KP_SCHEMA_REGISTRY_ENDPOINT")

		client, err := schemaregistry.NewClient(schemaregistry.NewConfig(os.Getenv("KP_SCHEMA_REGISTRY_ENDPOINT")))
		assert.NoError(t, err)
		ser, err := avro.NewGenericSerializer(client, serde.ValueSerde, avro.NewSerializerConfig())
		assert.NoError(t, err)
		for i := 0; i < 25000; i++ {
			payload, err := ser.Serialize("topic-kp", &BenchmarkMessage{
				Body:  "my-body",
				Count: 1000,
			})
			assert.NoError(t, err)
			kpPayload, err := serialization.Encode(BenchmarkMessage{
				Body:  "my-body",
				Count: 1000,
			}, 5)
			assert.NoError(t, err)
			assert.Equal(t, len(payload), len(kpPayload))
			assert.Equal(t, payload, kpPayload)
		}
	})
}
