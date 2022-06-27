package config

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/jinzhu/configor"
)

type Config struct {
	DB               DatabaseConfig
	Prometheus       PrometheusConfig
	Traces           TracesConfig
	KafkaConfig      KafkaConfig
	CoreAPI          CoreAPIConfig
	Anteraja         AnterajaConfig
	AppConfig        AppConfig
	GqlClientConfig  GqlClientConfig
	SchedulingConfig SchedulingConfig
	Shipper          ShipperConfig
}

type AppConfig struct {
	Name    string `env:"CONFIG__APP_CONFIG__NAME" required:"true" default:"card-delivery-service"`
	Version string `env:"APP__VERSION" default:"local"`
	Port    int    `env:"CONFIG__APP_CONFIG__PORT" default:"3000"`
}

type GqlClientConfig struct {
	Endpoint string `env:"CONFIG__GQL_CLIENT__ENDPOINT" default:"http://api-gateway-internal.honestcard.svc.cluster.local/graphql"`
}

type CoreAPIConfig struct {
	BaseURL string `env:"CONFIG__CORE_API_URL"`
}

type AnterajaConfig struct {
	BaseURL         string `env:"CONFIG__ANTERAJA_URL"`
	AccessKeyID     string `env:"CONFIG__ANTERAJA_ACCESS_KEY_ID"`
	SecretAccessKey string `env:"CONFIG__ANTERAJA_SECRET_ACCESS_KEY"`
	TimeZone        string `env:"CONFIG__TIME_ZONE"`
}

type SchedulingConfig struct {
	MaxCards int `env:"CONFIG__SCHEDULING__MAX_CARDS" default:"1400"`
}

type PrometheusConfig struct {
	PushGatewayHost string `env:"CONFIG__PROMETHEUS_CONFIG__HOST" default:"prometheus-pushgateway.monitoring.svc.cluster.local:9091"`
}

type TracesConfig struct {
	URL string `env:"CONFIG__TRACES_CONFIG__URL" default:"http://tempo.observability.svc.cluster.local:14268/api/traces"`
}

type DatabaseConfig struct {
	Host     string `env:"CONFIG__DB__HOST" required:"true"`
	Port     int    `env:"CONFIG__DB__PORT" default:"3306"`
	User     string `env:"CONFIG__DB__USER" required:"true"`
	Password string `env:"CONFIG__DB__PASSWORD" required:"true"`
	Name     string `env:"CONFIG__DB__NAME" required:"true"`
}

type KafkaConfig struct {
	KafkaBootstrapServers string `env:"CONFIG__KAFKA__BOOTSTRAP" default:"localhost:9092" json:"kafka_bootstrap"`
	KafkaTopicOutput      string `env:"CONFIG__KAFKA__TOPIC_OUTPUT" default:"card-emboss-request" json:"kafka_topic_output"`
	CardCreationTopic     string `env:"CONFIG__KAFKA__CARD_CREATION_TOPIC" required:"true"`
	CardReplacementTopic  string `env:"CONFIG__KAFKA__CARD_REPLACEMENT_TOPIC" required:"true"`
	CardRenewalTopic      string `env:"CONFIG__KAFKA__CARD_RENEWAL_TOPIC" required:"true"`
}

type ShipperConfig struct {
	Name       string `env:"CONFIG__SHIPPER__NAME" required:"true"`
	Phone      string `env:"CONFIG__SHIPPER__PHONE" required:"true"`
	Email      string `env:"CONFIG__SHIPPER__EMAIL" `
	District   string `env:"CONFIG__SHIPPER__DISTRICT" required:"true"`
	Address    string `env:"CONFIG__SHIPPER__ADDRESS" required:"true"`
	PostalCode string `env:"CONFIG__SHIPPER__POSTAL_CODE" required:"true"`
	GeoLoc     string `env:"CONFIG__SHIPPER__GEO_LOC"`
}

func LoadConfig() (*Config, error) {
	var config Config
	err := configor.
		New(&configor.Config{AutoReload: false}).
		Load(&config, fmt.Sprintf("%s/config.%s.json", getConfigLocation(), getEnv()))

	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c Config) GetKafkaBootstrap() []string {
	return strings.Split(c.KafkaConfig.KafkaBootstrapServers, ",")
}

func getConfigLocation() string {
	_, filename, _, _ := runtime.Caller(0)

	return path.Join(path.Dir(filename), "../config")
}

func getEnv() string {
	val := os.Getenv("APP_ENV")
	// todo: check our stage names and align with them
	switch val {
	case "prod":
		return "prod"
	case "test":
		return "test"
	case "qa":
		return "qa"
	default:
		return "dev"
	}
}
