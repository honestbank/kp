//go:build integration_test

package avro_serializer_test

import (
	"encoding/binary"
	"os"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/avro"
	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/internal/serialization/avro_serializer"
)

type BenchmarkMessage struct {
	Body  string `json:"body" avro:"body"`
	Count int    `json:"count" avro:"count"`
}

func TestSerialization(t *testing.T) {
	t.Run("can serialize and deserialize", func(t *testing.T) {
		serializer := avro_serializer.New[BenchmarkMessage](1)
		bytes, err := serializer.Encode(BenchmarkMessage{
			Body:  "my-body",
			Count: 100,
		})
		assert.NoError(t, err)
		assert.NotNil(t, bytes)
		msg, err := serializer.Decode(bytes)
		assert.NoError(t, err)
		assert.NotNil(t, msg)
		assert.Equal(t, "my-body", msg.Body)
		assert.Equal(t, 100, msg.Count)
	})
	t.Run("matches with confluent serializer", func(t *testing.T) {
		t.Setenv("KP_SCHEMA_REGISTRY_ENDPOINT", "http://localhost:8081")

		client, err := schemaregistry.NewClient(schemaregistry.NewConfig(os.Getenv("KP_SCHEMA_REGISTRY_ENDPOINT")))
		assert.NoError(t, err)
		ser, err := avro.NewGenericSerializer(client, serde.ValueSerde, avro.NewSerializerConfig())
		assert.NoError(t, err)
		for i := 0; i < 250; i++ {
			payload, err := ser.Serialize("topic-kp", &BenchmarkMessage{
				Body:  "my-body",
				Count: i,
			})
			assert.NoError(t, err)
			// do a hack to get schema id
			// schemaID := payload[1:5]
			schemaID := int(binary.BigEndian.Uint32(payload[1:5]))
			serializer := avro_serializer.New[BenchmarkMessage](schemaID)
			kpPayload, err := serializer.Encode(BenchmarkMessage{
				Body:  "my-body",
				Count: i,
			})
			assert.NoError(t, err)
			assert.Equal(t, len(payload), len(kpPayload))
			assert.Equal(t, payload, kpPayload)
		}
	})
}
