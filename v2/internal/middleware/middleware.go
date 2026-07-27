package middleware

import (
	"context"
	"errors"
	"io"
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
}

type stack[IN any, OUT any] struct {
	middlewares []Middleware[IN, OUT]
}

func (r *stack[IN, OUT]) AddMiddleware(mw Middleware[IN, OUT]) {
	r.middlewares = append(r.middlewares, mw)
}

func (r *stack[IN, OUT]) Process(ctx context.Context, options IN) OUT {
	var nextMiddleware func(c context.Context, item IN) OUT = nil
	middlewares := make([]Middleware[IN, OUT], len(r.middlewares))
	copy(middlewares, r.middlewares)
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
	for _, mw := range r.middlewares {
		if closer, ok := mw.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)
}

func New[IN, OUT any]() Processor[IN, OUT] {
	return &stack[IN, OUT]{}
}
