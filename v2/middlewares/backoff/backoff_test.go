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

// waitAfter waits for wait after 1 failure or more, and does not wait with no failures.
func waitAfter(wait time.Duration) func(failures int) time.Duration {
	return func(failures int) time.Duration {
		if failures == 0 {
			return 0
		}

		return wait
	}
}

func process(ctx context.Context, mw *backoff.InterruptibleBackoff, err error) (time.Duration, error) {
	start := time.Now()
	got := mw.Process(ctx, nil, func(ctx context.Context, msg *kafka.Message) error {
		return err
	})

	return time.Since(start), got
}

func TestInterruptibleBackoff(t *testing.T) {
	t.Run("does not wait before the first message and returns what next returns", func(t *testing.T) {
		mw := backoff.NewInterruptibleBackoffMiddleware(waitAfter(time.Hour))
		err := errors.New("some error")

		took, actualErr := process(context.Background(), mw, err)

		assert.Same(t, err, actualErr)
		assert.Less(t, took, time.Second)
	})

	t.Run("waits after a failure, and a success takes the failure away", func(t *testing.T) {
		mw := backoff.NewInterruptibleBackoffMiddleware(waitAfter(100 * time.Millisecond))
		_, _ = process(context.Background(), mw, errors.New("some error"))

		took, err := process(context.Background(), mw, nil)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, took, 100*time.Millisecond)

		took, _ = process(context.Background(), mw, nil)
		assert.Less(t, took, 100*time.Millisecond)
	})

	t.Run("Stop ends the wait in progress, and the message still runs", func(t *testing.T) {
		mw := backoff.NewInterruptibleBackoffMiddleware(waitAfter(time.Hour))
		_, _ = process(context.Background(), mw, errors.New("some error"))

		time.AfterFunc(50*time.Millisecond, mw.Stop)
		called := false
		start := time.Now()
		err := mw.Process(context.Background(), nil, func(ctx context.Context, msg *kafka.Message) error {
			called = true

			return nil
		})

		assert.NoError(t, err)
		assert.True(t, called)
		assert.Less(t, time.Since(start), 5*time.Second)
	})

	t.Run("a done context ends the wait", func(t *testing.T) {
		mw := backoff.NewInterruptibleBackoffMiddleware(waitAfter(time.Hour))
		_, _ = process(context.Background(), mw, errors.New("some error"))
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		took, _ := process(ctx, mw, nil)

		assert.Less(t, took, 5*time.Second)
	})

	t.Run("after Stop, the next messages do not wait, and a second Stop does not panic", func(t *testing.T) {
		mw := backoff.NewInterruptibleBackoffMiddleware(waitAfter(time.Hour))
		_, _ = process(context.Background(), mw, errors.New("some error"))
		mw.Stop()
		mw.Stop()

		took, _ := process(context.Background(), mw, errors.New("some error"))

		assert.Less(t, took, 5*time.Second)
	})
}
