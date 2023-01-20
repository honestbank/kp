package sigint_test

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/middlewares/sigint"
)

func TestSigInt_Process(t *testing.T) {
	t.Run("simply calls next", func(t *testing.T) {
		called := false
		sigint.NewSigIntMiddleware[*kafka.Message](func() {}).Process(context.Background(), nil, func(ctx context.Context, item *kafka.Message) error {
			called = true

			return nil
		})
		assert.True(t, called)
	})

	t.Run("when there's a sigint signal, it calls a callback", func(t *testing.T) {
		called := false
		sigint.NewSigIntMiddleware[*kafka.Message](func() {
			called = true
		}).Process(context.Background(), nil, func(ctx context.Context, item *kafka.Message) error {
			return nil
		})
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		time.Sleep(time.Millisecond * 50)
		assert.True(t, called)
	})
}
