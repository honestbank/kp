package sigint

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/honestbank/kp/v2/middlewares"
)

type sigInt[T any] struct{}

func (s sigInt[T]) Process(ctx context.Context, item T, next func(ctx context.Context, item T) error) error {
	return next(ctx, item)
}

func registerInterruptSignal(cb func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP)
	go func() {
		<-c
		cb()
	}()
}

func NewSigIntMiddleware[T any](stopFunction func()) middlewares.KPMiddleware[T] {
	registerInterruptSignal(stopFunction)

	return sigInt[T]{}
}
