package middleware

import (
	"context"
	"errors"
	"io"
	"sync"
)

type Middleware[IN any, OUT any] interface {
	Process(ctx context.Context, item IN, next func(ctx context.Context, item IN) OUT) OUT
}

type Processor[IN any, OUT any] interface {
	AddMiddleware(middleware Middleware[IN, OUT])
	Process(ctx context.Context, input IN) OUT
	// Close releases every middleware in the chain that implements io.Closer.
	// It lets the processor tear down resources (e.g. the Kafka consumer) once
	// processing has stopped, without the caller knowing which middlewares hold
	// them.
	Close() error
	// Stop stops every middleware in the chain that has a Stop method, such as the backoff, so that the message
	// in progress does not wait. It is safe to call from another goroutine while Process runs.
	Stop()
}

// stopper is a middleware that waits, and whose wait Stop ends.
type stopper interface {
	Stop()
}

type stack[IN any, OUT any] struct {
	// mu guards middlewares, because Stop can run in another goroutine while the chain is built or runs.
	mu          sync.RWMutex
	middlewares []Middleware[IN, OUT]
}

func (r *stack[IN, OUT]) AddMiddleware(mw Middleware[IN, OUT]) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.middlewares = append(r.middlewares, mw)
}

// snapshot returns a copy of the middlewares, so that the caller does not hold the lock while they run.
func (r *stack[IN, OUT]) snapshot() []Middleware[IN, OUT] {
	r.mu.RLock()
	defer r.mu.RUnlock()
	middlewares := make([]Middleware[IN, OUT], len(r.middlewares))
	copy(middlewares, r.middlewares)

	return middlewares
}

func (r *stack[IN, OUT]) Process(ctx context.Context, options IN) OUT {
	var nextMiddleware func(c context.Context, item IN) OUT = nil
	middlewares := r.snapshot()
	nextMiddleware = func(c context.Context, item IN) OUT {
		currentMw := middlewares[0]
		middlewares = middlewares[1:]
		return currentMw.Process(c, item, nextMiddleware)
	}
	return nextMiddleware(ctx, options)
}

// Close closes every middleware in the chain that implements io.Closer, in the
// order they were added, and joins any errors. Middlewares that are not closers
// are skipped.
func (r *stack[IN, OUT]) Close() error {
	var errs []error
	for _, mw := range r.snapshot() {
		if closer, ok := mw.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)
}

// Stop stops every middleware in the chain that has a Stop method, in the order they were added. The other
// middlewares are skipped.
func (r *stack[IN, OUT]) Stop() {
	for _, mw := range r.snapshot() {
		if s, ok := mw.(stopper); ok {
			s.Stop()
		}
	}
}

func New[IN, OUT any]() Processor[IN, OUT] {
	return &stack[IN, OUT]{}
}
