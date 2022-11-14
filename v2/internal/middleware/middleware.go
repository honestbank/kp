package middleware

import "context"

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

func New[IN, OUT any]() Processor[IN, OUT] {
	return &stack[IN, OUT]{}
}
