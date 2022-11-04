package middleware

import "context"

type finalMiddleware[IN, OUT any] struct {
	processor func(ctx context.Context, item IN) OUT
}

func (m finalMiddleware[IN, OUT]) Process(ctx context.Context, item IN, next func(ctx context.Context, item IN) OUT) OUT {
	return m.processor(ctx, item)
}

func FinalMiddleware[IN, OUT any](fn func(ctx context.Context, item IN) OUT) Middleware[IN, OUT] {
	return finalMiddleware[IN, OUT]{
		processor: fn,
	}
}

type Middleware[IN any, OUT any] interface {
	Process(ctx context.Context, item IN, next func(ctx context.Context, item IN) OUT) OUT
}

type Processor[IN any, OUT any] interface {
	AddMiddleware(middleware Middleware[IN, OUT])
	Process(ctx context.Context, input IN) OUT
}

type stack[IN any, OUT any] struct {
	middlewares []Middleware[IN, OUT]
}

func (r *stack[IN, OUT]) AddMiddleware(mw Middleware[IN, OUT]) {
	r.middlewares = append(r.middlewares, mw)
}

func (r *stack[IN, OUT]) Process(ctx context.Context, options IN) OUT {
	var next func(c context.Context, item IN) OUT = nil
	middlewares := append(r.middlewares)
	next = func(c context.Context, item IN) OUT {
		nextMw := middlewares[0]
		middlewares = middlewares[1:]
		return nextMw.Process(c, item, next)
	}
	return next(ctx, options)
}

func New[IN, OUT any]() Processor[IN, OUT] {
	return &stack[IN, OUT]{}
}
