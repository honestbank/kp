package noop_registry_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/internal/schemaregistry/noop_registry"
)

type test struct {
	Test string
}

func TestClient_Publish(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		client := noop_registry.New[test]()
		i, err := client.Publish()
		assert.Nil(t, i)
		assert.NoError(t, err)
	})
}
