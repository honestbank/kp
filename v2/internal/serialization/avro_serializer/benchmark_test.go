//go:build integration_test

package avro_serializer_test

import (
	"os"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/avro"
	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/internal/serialization/avro_serializer"
)

func BenchmarkEncode(b *testing.B) {
	os.Setenv("KP_SCHEMA_REGISTRY_ENDPOINT", "http://localhost:8081")
	defer os.Unsetenv("KP_SCHEMA_REGISTRY_ENDPOINT")
	client, err := schemaregistry.NewClient(schemaregistry.NewConfig(os.Getenv("KP_SCHEMA_REGISTRY_ENDPOINT")))
	assert.NoError(b, err)
	ser, err := avro.NewGenericSerializer(client, serde.ValueSerde, avro.NewSerializerConfig())
	assert.NoError(b, err)
	b.Run("kp", func(b *testing.B) {
		serializer := avro_serializer.New[BenchmarkMessage](1)
		for i := 0; i < b.N; i++ {
			_, err := serializer.Encode(BenchmarkMessage{
				Body:  "my-body",
				Count: 100,
			})
			assert.NoError(b, err)
		}
	})
	b.Run("kafka serializer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := ser.Serialize("topic-kp", &BenchmarkMessage{
				Body:  "my-body",
				Count: 1000,
			})
			assert.NoError(b, err)
		}
	})
}
