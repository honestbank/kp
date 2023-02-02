package backoff_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/honestbank/kp/v2/middlewares/backoff"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/assert"

	backoff_policy "github.com/honestbank/backoff-policy"
)

func TestBackoff(t *testing.T) {
	t.Run("calls next", func(t *testing.T) {
		called := false
		backoff.NewBackoffMiddleware(backoff_policy.NewExponentialBackoffPolicy(0, 0)).Process(context.Background(), nil, func(ctx context.Context, msg *kafka.Message) error {
			called = true
			return nil
		})
		assert.True(t, called)
	})
	t.Run("returns what next returns", func(t *testing.T) {
		mw := backoff.NewBackoffMiddleware(backoff_policy.NewExponentialBackoffPolicy(0, 0))
		err := errors.New("some error")
		actualErr := mw.Process(context.Background(), nil, func(ctx context.Context, msg *kafka.Message) error {
			return err
		})
		assert.Same(t, err, actualErr)
	})
	t.Run("when there's error, it slows down", func(t *testing.T) {
		mw := backoff.NewBackoffMiddleware(backoff_policy.NewExponentialBackoffPolicy(time.Second, 5))
		_ = mw.Process(context.Background(), nil, func(ctx context.Context, msg *kafka.Message) error {
			return errors.New("some error")
		})
		start := time.Now()
		_ = mw.Process(context.Background(), nil, func(ctx context.Context, msg *kafka.Message) error {
			return errors.New("some error")
		})
		assert.Greater(t, time.Since(start), time.Millisecond*1500)
		_ = mw.Process(context.Background(), nil, func(ctx context.Context, msg *kafka.Message) error {
			return errors.New("some error")
		})
		assert.Greater(t, time.Since(start), time.Millisecond*3500)
	})
}
