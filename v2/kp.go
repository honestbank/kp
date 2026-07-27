package v2

import (
	"context"
	"sync/atomic"

	"github.com/honestbank/kp/v2/internal/middleware"
)

type Processor[MessageType any] func(ctx context.Context, item *MessageType) error

type kp[MessageType any] struct {
	chain          middleware.Processor[*MessageType, error]
	shouldContinue *atomic.Bool
}

func (t *kp[MessageType]) getShouldContinue() bool {
	return t.shouldContinue.Load()
}

func (t *kp[MessageType]) AddMiddleware(middleware middleware.Middleware[*MessageType, error]) MessageProcessor[MessageType] {
	t.chain.AddMiddleware(middleware)

	return t
}

func (t *kp[MessageType]) Stop() {
	t.shouldContinue.Store(false)
}

func (t *kp[MessageType]) Run(processor Processor[MessageType]) error {
	t.chain.AddMiddleware(middleware.FinalMiddleware[*MessageType, error](func(ctx context.Context, msg *MessageType) error {
		return processor(ctx, msg)
	}))

	for t.getShouldContinue() {
		ctx := context.Background()
		_ = t.chain.Process(ctx, nil)
	}

	// Loop has stopped polling. Close the chain so resource-holding middlewares
	// (e.g. the Kafka consumer) release gracefully — the consumer leaves the
	// group here, letting the broker rebalance immediately. This runs
	// automatically, so downstream services no longer need a manual
	// defer kafkaConsumer.Close().
	return t.chain.Close()
}

func New[MessageType any]() MessageProcessor[MessageType] {
	return &kp[MessageType]{
		chain:          middleware.New[*MessageType, error](),
		shouldContinue: getAtomicBoolean(true),
	}
}

func getAtomicBoolean(value bool) *atomic.Bool {
	v := atomic.Bool{}
	v.Store(value)

	return &v
}
