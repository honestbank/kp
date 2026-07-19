package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/config"
)

func TestGetKafkaConfig(t *testing.T) {
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
	assignmentStrategy, err := withDefaults.Get("partition.assignment.strategy", "")
	assert.NoError(t, err)
	assert.Equal(t, "", assignmentStrategy, "must stay unset so librdkafka's default applies")
}

func TestGetKafkaConsumerConfig_PartitionAssignmentStrategy(t *testing.T) {
	strategy := "cooperative-sticky"
	cfg := config.GetKafkaConsumerConfig(config.Kafka{
		BootstrapServers:            "localhost",
		ConsumerGroupName:           "cg",
		PartitionAssignmentStrategy: &strategy,
	}.WithDefaults())
	got, err := cfg.Get("partition.assignment.strategy", "")
	assert.NoError(t, err)
	assert.Equal(t, "cooperative-sticky", got)
}
