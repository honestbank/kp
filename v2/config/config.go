package config

import "github.com/confluentinc/confluent-kafka-go/kafka"

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
}

func (s Kafka) WithDefaults() Kafka {
	return Kafka{
		ConsumerGroupName:        s.ConsumerGroupName,
		BootstrapServers:         s.BootstrapServers,
		SaslMechanism:            defaultIfNil(s.SaslMechanism, ""),
		SecurityProtocol:         defaultIfNil(s.SecurityProtocol, ""),
		Username:                 defaultIfNil(s.Username, ""),
		Password:                 defaultIfNil(s.Password, ""),
		ConsumerSessionTimeoutMs: defaultIfNil(s.ConsumerSessionTimeoutMs, 6000),
		ConsumerAutoOffsetReset:  defaultIfNil(s.ConsumerAutoOffsetReset, "earliest"),
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

	return cfg
}

func GetKafkaConsumerConfig(config Kafka) *kafka.ConfigMap {
	cfg := GetKafkaConfig(config)
	hydrateIfNotNil(cfg, "auto.offset.reset", config.ConsumerAutoOffsetReset)
	_ = cfg.SetKey("session.timeout.ms", config.ConsumerSessionTimeoutMs)

	return cfg
}

func hydrateIfNotNil(cfg *kafka.ConfigMap, key string, value *string) {
	if value == nil {
		return
	}
	// looked at the source code, as of now, there's no error being returned, it's always nil
	_ = cfg.SetKey(key, *value)
}
