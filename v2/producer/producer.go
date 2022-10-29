package producer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/honestbank/kp/v2/internal/config"

	"github.com/honestbank/kp/v2/internal/schemaregistry"
	"github.com/honestbank/kp/v2/internal/serialization"
)

type producer[BodyType any, KeyType KeyTypes] struct {
	schemaID int
	k        *kafka.Producer
	topic    string
}

func (p producer[BodyType, KeyType]) Produce(message KafkaMessage[BodyType, KeyType]) error {
	value, err := serialization.Encode(message.Body, p.schemaID)
	if err != nil {
		return err
	}

	partition := kafka.TopicPartition{
		Topic:     &p.topic,
		Partition: kafka.PartitionAny,
	}
	msg := &kafka.Message{
		TopicPartition: partition,
		Value:          value,
	}

	return p.k.Produce(msg, nil)
}

func (p producer[BodyType, KeyType]) Flush() error {
	// I just felt like 3000 because 3 second should be enough for messages to be flushed
	// we'll need to optimize this as we go
	p.k.Flush(3000)

	return nil
}

func New[MessageType any, KeyType KeyTypes](topic string) (Producer[MessageType, KeyType], error) {
	cfg, err := config.LoadConfig[config.KafkaConfig]()
	if err != nil {
		return nil, err
	}
	schemaID, err := schemaregistry.Publish[MessageType](topic)
	if err != nil {
		return nil, err
	}

	k, err := kafka.NewProducer(getKafkaConfig(*cfg))

	return producer[MessageType, KeyType]{
		k:        k,
		schemaID: *schemaID,
		topic:    topic,
	}, nil
}

func getKafkaConfig(kafkaConfig config.KafkaConfig) *kafka.ConfigMap {
	cfg := &kafka.ConfigMap{}

	hydrateIfNotNil(cfg, "bootstrap.servers", kafkaConfig.BootstrapServers)
	hydrateIfNotNil(cfg, "sasl.mechanisms", kafkaConfig.SaslMechanism)
	hydrateIfNotNil(cfg, "security.protocol", kafkaConfig.SecurityProtocol)
	hydrateIfNotNil(cfg, "sasl.username", kafkaConfig.Username)
	hydrateIfNotNil(cfg, "sasl.password", kafkaConfig.Password)

	return cfg
}

func hydrateIfNotNil(cfg *kafka.ConfigMap, key string, value *string) *kafka.ConfigMap {
	if value == nil {
		return cfg
	}
	// looked at the source code, as of now, there's no error being returned, it's always nil
	_ = cfg.SetKey(key, *value)

	return cfg
}
