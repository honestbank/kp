//go:build integration_test

package producer_test

import (
	"context"
	"os"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/avro"
	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/producer"
)

type BenchmarkMessage struct {
	Body  string `json:"body" avro:"body"`
	Count int    `json:"count" avro:"count"`
}

func BenchmarkProducer(b *testing.B) {
	os.Setenv("KP_SCHEMA_REGISTRY_ENDPOINT", "http://localhost:8081")
	defer os.Unsetenv("KP_SCHEMA_REGISTRY_ENDPOINT")

	kp, err := producer.New[BenchmarkMessage]("topic-kp")
	assert.NoError(b, err)
	defer kp.Flush()

	confluentProducer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	assert.NoError(b, err)
	client, err := schemaregistry.NewClient(schemaregistry.NewConfig(os.Getenv("KP_SCHEMA_REGISTRY_ENDPOINT")))
	assert.NoError(b, err)
	ser, err := avro.NewGenericSerializer(client, serde.ValueSerde, avro.NewSerializerConfig())
	assert.NoError(b, err)
	topic := "topic-confluent"
	defer confluentProducer.Flush(3000)

	b.Run("kp producer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := kp.Produce(context.Background(), BenchmarkMessage{
				Body:  "hello-world",
				Count: i,
			})
			assert.NoError(b, err)
		}
	})

	b.Run("confluent producer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			payload, err := ser.Serialize(topic, &BenchmarkMessage{
				Body:  "hello-world",
				Count: i,
			})
			assert.NoError(b, err)
			err = confluentProducer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          payload,
			}, nil)
			assert.NoError(b, err)
		}
	})
}
