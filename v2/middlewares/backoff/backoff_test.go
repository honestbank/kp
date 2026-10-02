package backoff_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/honestbank/kp/v2/middlewares"
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
func waitAfter(wait time.Duration) func(count int) time.Duration {
	return func(count int) time.Duration {
		if count == 0 {
			return 0
		}

		return wait
	}
}

// executeOnly is a policy with Execute only, as a policy from before ContextBackoffPolicy.
type executeOnly struct {
	p backoff_policy.BackoffPolicy
}

func (e executeOnly) Execute(cb func(marker backoff_policy.Marker)) {
	e.p.Execute(cb)
}

func process(ctx context.Context, mw middlewares.KPMiddleware[*kafka.Message], err error) (time.Duration, error) {
	start := time.Now()
	got := mw.Process(ctx, nil, func(ctx context.Context, msg *kafka.Message) error {
		return err
	})

	return time.Since(start), got
}

func stop(t *testing.T, mw middlewares.KPMiddleware[*kafka.Message]) func() {
	t.Helper()
	stopper, ok := mw.(interface{ Stop() })
	assert.True(t, ok, "the backoff middleware has Stop")

	return stopper.Stop
}

func TestBackoffStop(t *testing.T) {
	t.Run("Stop ends the wait in progress, and the message still runs with its context", func(t *testing.T) {
		mw := backoff.NewBackoffMiddleware(backoff_policy.NewBackoff(waitAfter(time.Hour)))
		_, _ = process(context.Background(), mw, errors.New("some error"))

		time.AfterFunc(50*time.Millisecond, stop(t, mw))
		var processCtx context.Context
		start := time.Now()
		err := mw.Process(context.Background(), nil, func(ctx context.Context, msg *kafka.Message) error {
			processCtx = ctx

			return nil
		})

		assert.NoError(t, err)
		assert.Less(t, time.Since(start), 5*time.Second)
		assert.NoError(t, processCtx.Err(), "next gets the context of the message, not the ended wait context")
	})

	t.Run("a done context ends the wait", func(t *testing.T) {
		mw := backoff.NewBackoffMiddleware(backoff_policy.NewBackoff(waitAfter(time.Hour)))
		_, _ = process(context.Background(), mw, errors.New("some error"))
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		took, _ := process(ctx, mw, nil)

		assert.Less(t, took, 5*time.Second)
	})

	t.Run("after Stop, the next messages do not wait, and a second Stop does not panic", func(t *testing.T) {
		mw := backoff.NewBackoffMiddleware(backoff_policy.NewBackoff(waitAfter(time.Hour)))
		_, _ = process(context.Background(), mw, errors.New("some error"))
		stop(t, mw)()
		stop(t, mw)()

		took, _ := process(context.Background(), mw, errors.New("some error"))

		assert.Less(t, took, 5*time.Second)
	})

	t.Run("a policy with Execute only still waits, and Stop does not end its wait", func(t *testing.T) {
		mw := backoff.NewBackoffMiddleware(executeOnly{backoff_policy.NewBackoff(waitAfter(100 * time.Millisecond))})
		_, _ = process(context.Background(), mw, errors.New("some error"))
		stop(t, mw)()

		took, _ := process(context.Background(), mw, nil)

		assert.GreaterOrEqual(t, took, 100*time.Millisecond)
	})
}
