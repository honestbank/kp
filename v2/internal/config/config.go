package config

import "github.com/jinzhu/configor"

type KafkaConfig struct {
	BootstrapServers *string `env:"KP_KAFKA_BOOTSTRAP_SERVERS" json:"bootstrap_servers" default:""`
	SaslMechanism    *string `env:"KP_SASL_MECHANISM" json:"sasl_mechanism" default:""`
	SecurityProtocol *string `env:"KP_SECURITY_PROTOCOL" json:"security_protocol" default:""`
	Username         *string `env:"KP_USERNAME" json:"username" default:""`
	Password         *string `env:"KP_PASSWORD" json:"password" default:""`
}

func LoadConfig[T any]() (*T, error) {
	var cfg T
	err := configor.Load(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
