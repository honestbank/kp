package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/config"
)

func TestGetKafkaConfig(t *testing.T) {
	// WithDefaults now defaults to the KIP-848 consumer protocol.
	withDefaults := config.GetKafkaConsumerConfig(config.Kafka{BootstrapServers: "localhost", ConsumerGroupName: "cg"}.WithDefaults())
	bootstrapServers, err := withDefaults.Get("bootstrap.servers", "")
	assert.NoError(t, err)
	assert.Equal(t, "localhost", bootstrapServers)
	autoOffsetResets, err := withDefaults.Get("auto.offset.reset", "")
	assert.NoError(t, err)
	assert.Equal(t, "earliest", autoOffsetResets)
	consumerGroupId, err := withDefaults.Get("group.id", "")
	assert.NoError(t, err)
	assert.Equal(t, "cg", consumerGroupId)
	groupProtocol, err := withDefaults.Get("group.protocol", "")
	assert.NoError(t, err)
	assert.Equal(t, "consumer", groupProtocol, "defaults to the KIP-848 consumer protocol")
	// Under the consumer protocol these are server-managed and must not be sent client-side.
	sessionTimeout, err := withDefaults.Get("session.timeout.ms", "unset")
	assert.NoError(t, err)
	assert.Equal(t, "unset", sessionTimeout, "must not be set under the consumer protocol")
	assignmentStrategy, err := withDefaults.Get("partition.assignment.strategy", "unset")
	assert.NoError(t, err)
	assert.Equal(t, "unset", assignmentStrategy, "must not be set under the consumer protocol")
}

func TestGetKafkaConsumerConfig_ClassicProtocol(t *testing.T) {
	protocol := "classic"
	cfg := config.GetKafkaConsumerConfig(config.Kafka{
		BootstrapServers:  "localhost",
		ConsumerGroupName: "cg",
		GroupProtocol:     &protocol,
	}.WithDefaults())

	got, err := cfg.Get("group.protocol", "")
	assert.NoError(t, err)
	assert.Equal(t, "classic", got)
	// Classic protocol keeps the client-side session timeout default.
	sessionTimeout, err := cfg.Get("session.timeout.ms", 0)
	assert.NoError(t, err)
	assert.Equal(t, 30000, sessionTimeout, "classic protocol keeps the client-side default")
	assignmentStrategy, err := cfg.Get("partition.assignment.strategy", "")
	assert.NoError(t, err)
	assert.Equal(t, "", assignmentStrategy, "must stay unset so librdkafka's default applies")
}

func TestGetKafkaConsumerConfig_ConsumerProtocol(t *testing.T) {
	protocol := "consumer"
	strategy := "cooperative-sticky"
	cfg := config.GetKafkaConsumerConfig(config.Kafka{
		BootstrapServers:            "localhost",
		ConsumerGroupName:           "cg",
		GroupProtocol:               &protocol,
		PartitionAssignmentStrategy: &strategy,
	}.WithDefaults())

	got, err := cfg.Get("group.protocol", "")
	assert.NoError(t, err)
	assert.Equal(t, "consumer", got)

	// KIP-848: these are server-managed and must not be sent client-side, even though
	// WithDefaults() populates a session timeout and the caller passed a strategy.
	sessionTimeout, err := cfg.Get("session.timeout.ms", "unset")
	assert.NoError(t, err)
	assert.Equal(t, "unset", sessionTimeout, "must not be set under the consumer protocol")
	assignmentStrategy, err := cfg.Get("partition.assignment.strategy", "unset")
	assert.NoError(t, err)
	assert.Equal(t, "unset", assignmentStrategy, "must not be set under the consumer protocol")
}

func TestGetKafkaConsumerConfig_PartitionAssignmentStrategy(t *testing.T) {
	// PartitionAssignmentStrategy only applies under the classic protocol, so opt out
	// of the consumer-protocol default explicitly.
	protocol := "classic"
	strategy := "cooperative-sticky"
	cfg := config.GetKafkaConsumerConfig(config.Kafka{
		BootstrapServers:            "localhost",
		ConsumerGroupName:           "cg",
		GroupProtocol:               &protocol,
		PartitionAssignmentStrategy: &strategy,
	}.WithDefaults())
	got, err := cfg.Get("partition.assignment.strategy", "")
	assert.NoError(t, err)
	assert.Equal(t, "cooperative-sticky", got)
}
