package middleware_test

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/internal/middleware"
)

// stopperMw is a middleware with a Stop method, which records whether Stop was called.
type stopperMw struct {
	stopped *atomic.Bool
}

func (s stopperMw) Process(ctx context.Context, item int, next func(ctx context.Context, item int) int) int {
	return next(ctx, item)
}

func (s stopperMw) Stop() {
	s.stopped.Store(true)
}

func TestStackStop(t *testing.T) {
	t.Run("stops the middlewares that have Stop and skips the others", func(t *testing.T) {
		var stoppedA, stoppedB atomic.Bool
		stack := middleware.New[int, int]()
		stack.AddMiddleware(stopperMw{stopped: &stoppedA})
		stack.AddMiddleware(plainMw{}) // must be skipped, not panic
		stack.AddMiddleware(stopperMw{stopped: &stoppedB})

		stack.Stop()

		assert.True(t, stoppedA.Load(), "first stopper should be stopped")
		assert.True(t, stoppedB.Load(), "second stopper should be stopped")
	})

	t.Run("is safe while the chain is built and runs in another goroutine", func(t *testing.T) {
		var stopped atomic.Bool
		stack := middleware.New[int, int]()
		done := make(chan struct{})
		go func() {
			defer close(done)
			for range 100 {
				stack.Stop()
			}
		}()
		for range 100 {
			stack.AddMiddleware(stopperMw{stopped: &stopped})
		}
		stack.AddMiddleware(middleware.FinalMiddleware[int, int](func(ctx context.Context, item int) int { return item }))
		for range 100 {
			assert.Equal(t, 1, stack.Process(context.Background(), 1))
		}
		<-done
	})
}
