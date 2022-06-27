package kp

type Config struct {
	KafkaConfig KafkaConfig `json:"kafka"`
}

type KafkaConfig struct {
	KafkaBootstrapServers string `json:"kafka_bootstrap"`
}
