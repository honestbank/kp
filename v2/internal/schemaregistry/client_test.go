//go:build integration_test

package schemaregistry_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/internal/schemaregistry"
)

type BenchmarkMessage struct {
	Body  string `json:"body" avro:"body"`
	Count int    `json:"count" avro:"count"`
}

func TestPublish(t *testing.T) {
	os.Setenv("KP_SCHEMA_REGISTRY_ENDPOINT", "http://localhost:8081")
	defer os.Unsetenv("KP_SCHEMA_REGISTRY_ENDPOINT")
	_, err := schemaregistry.Publish[BenchmarkMessage]("topic-kp")
	assert.NoError(t, err)
}
