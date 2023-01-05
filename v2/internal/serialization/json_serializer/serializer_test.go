package json_serializer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/internal/serialization/json_serializer"
)

type message struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

func TestSerializer_Decode(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		serializer := json_serializer.New[message]()
		m, err := serializer.Decode(nil)
		assert.Error(t, err)
		assert.Nil(t, m)
	})
	t.Run("happy", func(t *testing.T) {
		serializer := json_serializer.New[message]()
		payload := []byte(`{"id": "123", "value": "test"}`)
		m, err := serializer.Decode(payload)
		assert.NoError(t, err)
		assert.Equal(t, "123", m.ID)
		assert.Equal(t, "test", m.Value)
	})
}

func TestSerializer_Encode(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		serializer := json_serializer.New[message]()
		m := message{ID: "123", Value: "test"}
		b, err := serializer.Encode(m)
		assert.NoError(t, err)
		assert.Contains(t, string(b), "123")
	})
}
