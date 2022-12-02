package v2

import (
	"context"

	"github.com/honestbank/kp/v2/internal/middleware"
)

type Processor[MessageType any] func(ctx context.Context, item *MessageType) error

type kp[MessageType any] struct {
	chain          middleware.Processor[*MessageType, error]
	shouldContinue bool
}

func (t *kp[MessageType]) AddMiddleware(middleware middleware.Middleware[*MessageType, error]) KafkaProcessor[MessageType] {
	t.chain.AddMiddleware(middleware)

	return t
}

func (t *kp[MessageType]) Stop() {
	t.shouldContinue = false
}

func (t *kp[MessageType]) Run(processor Processor[MessageType]) error {
	t.chain.AddMiddleware(middleware.FinalMiddleware[*MessageType, error](func(ctx context.Context, msg *MessageType) error {
		return processor(ctx, msg)
	}))
	for t.shouldContinue {
		ctx := context.Background()
		_ = t.chain.Process(ctx, nil)
	}
	return nil
}

func New[MessageType any]() KafkaProcessor[MessageType] {
	return &kp[MessageType]{
		chain:          middleware.New[*MessageType, error](),
		shouldContinue: true,
	}
}
