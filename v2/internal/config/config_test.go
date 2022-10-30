package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/internal/config"
)

type Cfg struct {
	SomeValue string `json:"some_value" env:"SOME_VALUE" required:"true"`
}

func TestLoadConfig(t *testing.T) {
	t.Setenv("KP_KAFKA_BOOTSTRAP_SERVERS", "localhost")
	cfg, err := config.LoadConfig[config.KafkaConfig]()
	assert.NoError(t, err)
	assert.Nil(t, cfg.Username)
	assert.Nil(t, cfg.Password)
	assert.Equal(t, "localhost", *cfg.BootstrapServers)
	assert.Nil(t, cfg.SecurityProtocol)
	assert.Nil(t, cfg.SaslMechanism)
}

func TestLoadConfigReturnsError(t *testing.T) {
	cfg, err := config.LoadConfig[Cfg]()
	assert.Error(t, err)
	assert.Nil(t, cfg)
}
