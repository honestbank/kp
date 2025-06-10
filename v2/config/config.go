package config

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type KPConfig struct {
	KafkaConfig          Kafka
	SchemaRegistryConfig SchemaRegistry
}

type Kafka struct {
	ConsumerGroupName        string
	BootstrapServers         string
	SaslMechanism            *string
	SecurityProtocol         *string
	Username                 *string
	Password                 *string
	ConsumerSessionTimeoutMs *int
	ConsumerAutoOffsetReset  *string
	ClientID                 *string
	MaxMessageBytes          *int
	Debug                    *string
}

func (s Kafka) WithDefaults() Kafka {
	return Kafka{
		ConsumerGroupName:        s.ConsumerGroupName,
		BootstrapServers:         s.BootstrapServers,
		SaslMechanism:            s.SaslMechanism,
		SecurityProtocol:         s.SecurityProtocol,
		Username:                 s.Username,
		Password:                 s.Password,
		ConsumerSessionTimeoutMs: defaultIfNil(s.ConsumerSessionTimeoutMs, 30000),
		ConsumerAutoOffsetReset:  defaultIfNil(s.ConsumerAutoOffsetReset, "earliest"),
		ClientID:                 defaultIfNil(s.ClientID, "rdkafka"),
		Debug:                    s.Debug,
	}
}

type SchemaRegistry struct {
	Endpoint string
	Username string
	Password string
}

func defaultIfNil[T any](value *T, defaultValue T) *T {
	if value == nil {
		return &defaultValue
	}

	return value
}

func GetKafkaConfig(kafkaConfig Kafka) *kafka.ConfigMap {
	cfg := &kafka.ConfigMap{}

	hydrateIfNotNil(cfg, "bootstrap.servers", &kafkaConfig.BootstrapServers)
	hydrateIfNotNil(cfg, "sasl.mechanisms", kafkaConfig.SaslMechanism)
	hydrateIfNotNil(cfg, "security.protocol", kafkaConfig.SecurityProtocol)
	hydrateIfNotNil(cfg, "sasl.username", kafkaConfig.Username)
	hydrateIfNotNil(cfg, "sasl.password", kafkaConfig.Password)
	hydrateIfNotNil(cfg, "debug", kafkaConfig.Debug)
	hydrateIfNotNil(cfg, "client.id", kafkaConfig.ClientID)
	hydrateIfNotNil(cfg, "message.max.bytes", kafkaConfig.MaxMessageBytes)

	return cfg
}

func GetKafkaConsumerConfig(config Kafka) *kafka.ConfigMap {
	cfg := GetKafkaConfig(config)
	hydrateIfNotNil(cfg, "group.id", &config.ConsumerGroupName)
	hydrateIfNotNil(cfg, "auto.offset.reset", config.ConsumerAutoOffsetReset)
	hydrateIfNotNil(cfg, "session.timeout.ms", config.ConsumerSessionTimeoutMs)
	hydrateIfNotNil(cfg, "debug", config.Debug)

	return cfg
}

func hydrateIfNotNil[T any](cfg *kafka.ConfigMap, key string, value *T) {
	if value == nil {
		return
	}
	// looked at the source code, as of now, there's no error being returned, it's always nil
	_ = cfg.SetKey(key, *value)
}
