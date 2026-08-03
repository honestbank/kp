package middleware_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/internal/middleware"
)

// closerMw is a middleware that also implements io.Closer, recording whether
// Close was called and optionally returning an error.
type closerMw struct {
	closed *bool
	err    error
}

func (c closerMw) Process(ctx context.Context, item int, next func(ctx context.Context, item int) int) int {
	return next(ctx, item)
}

func (c closerMw) Close() error {
	*c.closed = true
	return c.err
}

// plainMw is a middleware that does NOT implement io.Closer.
type plainMw struct{}

func (plainMw) Process(ctx context.Context, item int, next func(ctx context.Context, item int) int) int {
	return next(ctx, item)
}

func TestStackClose(t *testing.T) {
	t.Run("closes io.Closer middlewares and skips non-closers", func(t *testing.T) {
		closedA, closedB := false, false
		stack := middleware.New[int, int]()
		stack.AddMiddleware(closerMw{closed: &closedA})
		stack.AddMiddleware(plainMw{}) // must be skipped, not panic
		stack.AddMiddleware(closerMw{closed: &closedB})

		err := stack.Close()

		assert.NoError(t, err)
		assert.True(t, closedA, "first closer should be closed")
		assert.True(t, closedB, "second closer should be closed")
	})

	t.Run("joins errors from failing closers", func(t *testing.T) {
		errA := errors.New("close A failed")
		errB := errors.New("close B failed")
		closedA, closedB := false, false
		stack := middleware.New[int, int]()
		stack.AddMiddleware(closerMw{closed: &closedA, err: errA})
		stack.AddMiddleware(closerMw{closed: &closedB, err: errB})

		err := stack.Close()

		assert.ErrorIs(t, err, errA)
		assert.ErrorIs(t, err, errB)
		assert.True(t, closedA)
		assert.True(t, closedB)
	})

	t.Run("no closers is a no-op", func(t *testing.T) {
		stack := middleware.New[int, int]()
		stack.AddMiddleware(plainMw{})

		assert.NoError(t, stack.Close())
	})
}
