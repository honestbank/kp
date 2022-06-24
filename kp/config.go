package kp

type Config struct {
	KafkaConfig KafkaConfig `json:"kafka"`
}

type KafkaConfig struct {
	KafkaBootstrap string `env:"CONFIG__KAFKA__BOOTSTRAP" default:"localhost:9092" json:"kafka_bootstrap"`
}
