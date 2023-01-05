//go:build integration_test

package avro_registry_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/config"
	"github.com/honestbank/kp/v2/internal/schemaregistry/avro_registry"
)

type BenchmarkMessage struct {
	Body  string `json:"body" avro:"body"`
	Count int    `json:"count" avro:"count"`
}

func TestPublish(t *testing.T) {
	client, err := avro_registry.New[BenchmarkMessage]("topic-kp", config.SchemaRegistry{Endpoint: "http://localhost:8081"})
	assert.NoError(t, err)

	_, err = client.Publish()
	assert.NoError(t, err)
}
