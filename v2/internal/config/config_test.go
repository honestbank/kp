package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/internal/config"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := config.LoadConfig[config.KafkaConfig]()
	assert.NoError(t, err)
	assert.Nil(t, cfg.Username)
	assert.Nil(t, cfg.Password)
	assert.Nil(t, cfg.BootstrapServers)
	assert.Nil(t, cfg.SecurityProtocol)
	assert.Nil(t, cfg.SaslMechanism)
}
