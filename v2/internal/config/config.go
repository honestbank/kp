package config

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jinzhu/configor"
)

type KafkaConfig struct {
	BootstrapServers         *string `env:"KP_KAFKA_BOOTSTRAP_SERVERS" json:"bootstrap_servers" required:"true"`
	SaslMechanism            *string `env:"KP_SASL_MECHANISM" json:"sasl_mechanism" default:""`
	SecurityProtocol         *string `env:"KP_SECURITY_PROTOCOL" json:"security_protocol" default:""`
	Username                 *string `env:"KP_USERNAME" json:"username" default:""`
	Password                 *string `env:"KP_PASSWORD" json:"password" default:""`
	ConsumerSessionTimeoutMs int     `env:"KP_SESSION_TIMEOUT_MS" json:"consumer_session_timeout_ms" default:"6000"`
	ConsumerAutoOffsetReset  *string `env:"KP_CONSUMER_AUTO_OFFSET_RESET" json:"consumer_auto_offset_reset" default:"earliest"`
}

func LoadConfig[T any]() (*T, error) {
	var cfg T
	err := configor.Load(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func GetKafkaConfig(kafkaConfig KafkaConfig) *kafka.ConfigMap {
	cfg := &kafka.ConfigMap{}

	hydrateIfNotNil(cfg, "bootstrap.servers", kafkaConfig.BootstrapServers)
	hydrateIfNotNil(cfg, "sasl.mechanisms", kafkaConfig.SaslMechanism)
	hydrateIfNotNil(cfg, "security.protocol", kafkaConfig.SecurityProtocol)
	hydrateIfNotNil(cfg, "sasl.username", kafkaConfig.Username)
	hydrateIfNotNil(cfg, "sasl.password", kafkaConfig.Password)

	return cfg
}

func GetKafkaConsumerConfig(config KafkaConfig) *kafka.ConfigMap {
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
