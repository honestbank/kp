package config

import "github.com/jinzhu/configor"

type KafkaConfig struct {
	BootstrapServers string `json:"bootstrap_servers" default:""`
}

func LoadConfig[T any]() (*T, error) {
	var cfg T
	err := configor.Load(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
