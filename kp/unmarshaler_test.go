package kp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/kp"
)

func TestUnmarshalStringMessage(t *testing.T) {
	t.Run("test unmarshal string message", func(t *testing.T) {
		a := assert.New(t)

		message := "test|1"
		key, retries, err := kp.UnmarshalStringMessage(message)
		a.NoError(err)
		a.Equal("test", key)
		a.Equal(1, retries)
	})
	t.Run("test unmarshal string message with no retries", func(t *testing.T) {
		a := assert.New(t)

		message := "test"
		key, retries, err := kp.UnmarshalStringMessage(message)
		a.NoError(err)
		a.Equal("test", key)
		a.Equal(0, retries)
	})

	t.Run("test unmarshal string message with \"|\" and no number", func(t *testing.T) {
		a := assert.New(t)

		message := "test|"
		key, retries, _ := kp.UnmarshalStringMessage(message)
		a.Equal("test|", key)
		a.Equal(0, retries)
	})

	t.Run("test unmarshal string message multiple \"|\"", func(t *testing.T) {
		a := assert.New(t)

		message := "test|1|2"
		key, retries, err := kp.UnmarshalStringMessage(message)
		a.NoError(err)
		a.Equal("test|1", key)
		a.Equal(2, retries)
	})

	t.Run("test unmarshal string message with no retries", func(t *testing.T) {
		a := assert.New(t)

		message := "test"
		msg, retries, err := kp.UnmarshalStringMessage(message)
		a.NoError(err)
		a.Equal("test", msg)
		a.Equal(0, retries)
	})
}

func TestMarshalStringMessage(t *testing.T) {
	t.Run("test marshal string message", func(t *testing.T) {
		a := assert.New(t)

		message := "test"
		msg := kp.MarshalStringMessage(message, 1)
		a.Equal("test|1", msg)
	})
	t.Run("test marshal string message - no retries", func(t *testing.T) {
		a := assert.New(t)

		message := "test"
		msg := kp.MarshalStringMessage(message, 0)
		a.Equal("test", msg)
	})
}
