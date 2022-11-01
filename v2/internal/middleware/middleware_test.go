package middleware_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/kp/v2/internal/middleware"
)

type customMw struct {
}

func (r customMw) Process(item int, next func(item int) int) int {
	return next(item + 5)
}

func TestNew(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		pipeline := middleware.New[int, int]()
		pipeline.AddMiddleware(middleware.FinalMiddleware[int, int](func(item int) int { return item * item }))
		assert.Equal(t, 1, pipeline.Process(1))
		assert.Equal(t, 25, pipeline.Process(5))
	})
	t.Run("with middleware", func(t *testing.T) {
		pipeline := middleware.New[int, int]()
		pipeline.AddMiddleware(customMw{})
		pipeline.AddMiddleware(middleware.FinalMiddleware[int, int](func(item int) int { return item * item }))
		assert.Equal(t, 36, pipeline.Process(1))
		assert.Equal(t, 100, pipeline.Process(5))
	})
}
