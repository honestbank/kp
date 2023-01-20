package sigint

import (
	"context"
	"os"
	"os/signal"

	"github.com/honestbank/kp/v2/middlewares"
)

type sigInt[T any] struct{}

func (s sigInt[T]) Process(ctx context.Context, item T, next func(ctx context.Context, item T) error) error {
	return next(ctx, item)
}

func registerInterruptSignal(cb func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cb()
	}()
}

func NewSigIntMiddleware[T any](stopFunction func()) middlewares.KPMiddleware[T] {
	registerInterruptSignal(stopFunction)

	return sigInt[T]{}
}
