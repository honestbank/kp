package middleware

type finalMw[IN, OUT any] struct {
	processor func(item IN) OUT
}

func (m finalMw[IN, OUT]) Process(item IN, next func(item IN) OUT) OUT {
	return m.processor(item)
}

func FinalMw[IN, OUT any](fn func(item IN) OUT) Middleware[IN, OUT] {
	return finalMw[IN, OUT]{
		processor: fn,
	}
}

type Middleware[IN any, OUT any] interface {
	Process(item IN, next func(item IN) OUT) OUT
}

type Processor[IN any, OUT any] interface {
	AddMw(middleware Middleware[IN, OUT])
	Process(input IN) OUT
}

type stack[IN any, OUT any] struct {
	middlewares []Middleware[IN, OUT]
}

func (r *stack[IN, OUT]) AddMw(mw Middleware[IN, OUT]) {
	r.middlewares = append(r.middlewares, mw)
}

func (r *stack[IN, OUT]) Process(options IN) OUT {
	var next func(item IN) OUT = nil
	middlewares := append(r.middlewares)
	next = func(item IN) OUT {
		nextMw := middlewares[0]
		middlewares = middlewares[1:]
		return nextMw.Process(item, next)
	}
	return next(options)
}

func New[IN, OUT any]() Processor[IN, OUT] {
	return &stack[IN, OUT]{}
}
