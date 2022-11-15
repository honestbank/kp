//go:build integration_test

package schemaregistry_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/config"
	"github.com/honestbank/kp/v2/internal/schemaregistry"
)

type BenchmarkMessage struct {
	Body  string `json:"body" avro:"body"`
	Count int    `json:"count" avro:"count"`
}

func TestPublish(t *testing.T) {
	_, err := schemaregistry.Publish[BenchmarkMessage]("topic-kp", config.SchemaRegistry{Endpoint: "http://localhost:8081"})
	assert.NoError(t, err)
}
