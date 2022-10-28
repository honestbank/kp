//go:build integration_test

package producer_test

import (
	"os"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/avro"
	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/internal/serialization"
	"github.com/honestbank/kp/v2/producer"
)

type MyMessage struct {
	Id    string `json:"id" avro:"id"`
	Count int    `json:"count" avro:"count"`
}

type MyMessageBreaking struct {
	Id     string `json:"id" avro:"id"`
	Count2 int    `json:"count2" avro:"count2"`
}

func TestNewProducer(t *testing.T) {
	os.Setenv("KP_SCHEMA_REGISTRY_ENDPOINT", "http://localhost:8081")
	os.Setenv("KP_KAFKA_BOOTSTRAP_SERVERS", "localhost")
	defer os.Unsetenv("KP_SCHEMA_REGISTRY_ENDPOINT")
	defer os.Unsetenv("KP_KAFKA_BOOTSTRAP_SERVERS")
	t.Run("schema registry", func(t *testing.T) {
		t.Run("when a producer is initialized, schema is automatically registered", func(t *testing.T) {
			_, err := producer.New[MyMessage, int]("test-topic-3")
			assert.NoError(t, err)
			client, err := schemaregistry.NewClient(schemaregistry.NewConfigWithAuthentication(
				"http://localhost:8081",
				"",
				"",
			))
			assert.NoError(t, err)
			res, err := client.GetLatestSchemaMetadata("test-topic-3-value")
			assert.NoError(t, err)
			assert.Equal(t, "test-topic-3-value", res.Subject)
		})
		t.Run("fails initializing producer if there's a breaking change in schema", func(t *testing.T) {
			_, err := producer.New[MyMessage, int]("test-topic-1")
			_, err = producer.New[MyMessageBreaking, int]("test-topic-1")
			assert.Error(t, err)
		})
		t.Run("non breaking change allows initialization", func(t *testing.T) {
			_, err := producer.New[MyMessage, int]("test-topic-2")
			_, err = producer.New[MyMessage, int]("test-topic-2")
			assert.NoError(t, err)
		})
	})
}

func TestNew(t *testing.T) {
	os.Setenv("KP_SCHEMA_REGISTRY_ENDPOINT", "http://localhost:8081")
	defer os.Unsetenv("KP_SCHEMA_REGISTRY_ENDPOINT")
	t.Run("produce through confluent", func(t *testing.T) {
		confluentProducer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
		assert.NoError(t, err)
		client, err := schemaregistry.NewClient(schemaregistry.NewConfig(os.Getenv("KP_SCHEMA_REGISTRY_ENDPOINT")))
		assert.NoError(t, err)
		ser, err := avro.NewGenericSerializer(client, serde.ValueSerde, avro.NewSerializerConfig())
		assert.NoError(t, err)
		topic := "topic-confluent"
		defer confluentProducer.Flush(3000)

		for i := 0; i < 25000; i++ {
			payload, err := ser.Serialize(topic, &BenchmarkMessage{
				Body:  "hello-world",
				Count: i,
			})
			assert.NoError(t, err)
			kPayload, err := serialization.Encode(BenchmarkMessage{
				Body:  "hello-world",
				Count: i,
			}, 5)
			assert.Equal(t, kPayload, payload)
			assert.Equal(t, len(kPayload), len(payload))
			err = confluentProducer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          payload,
			}, nil)
			assert.NoError(t, err)
		}
	})
	t.Run("produce through kp", func(t *testing.T) {
		kp, err := producer.New[BenchmarkMessage, int]("topic-kp")
		assert.NoError(t, err)
		defer kp.Flush()

		for i := 0; i < 25000; i++ {
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
