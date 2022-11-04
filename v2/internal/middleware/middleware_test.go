package middleware_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/internal/middleware"
)

type customMw struct {
}

func (r customMw) Process(ctx context.Context, item int, next func(ctx context.Context, item int) int) int {
	return next(context.WithValue(ctx, "key", "some-value"), item+5)
}

func TestNew(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		pipeline := middleware.New[int, int]()
		pipeline.AddMiddleware(middleware.FinalMiddleware[int, int](func(ctx context.Context, item int) int { return item * item }))
		assert.Equal(t, 1, pipeline.Process(context.Background(), 1))
		assert.Equal(t, 25, pipeline.Process(context.Background(), 5))
	})
	t.Run("with middleware", func(t *testing.T) {
		pipeline := middleware.New[int, int]()
		pipeline.AddMiddleware(customMw{})
		pipeline.AddMiddleware(middleware.FinalMiddleware[int, int](func(ctx context.Context, item int) int { return item * item }))
		assert.Equal(t, 36, pipeline.Process(context.Background(), 1))
		assert.Equal(t, 100, pipeline.Process(context.Background(), 5))
	})

	t.Run("with context values", func(t *testing.T) {
		pipeline := middleware.New[int, int]()
		pipeline.AddMiddleware(customMw{})
		pipeline.AddMiddleware(middleware.FinalMiddleware[int, int](func(ctx context.Context, item int) int {
			assert.Equal(t, "some-value", ctx.Value("key"))
			return item
		}))
		pipeline.Process(context.Background(), 20)
	})
}
