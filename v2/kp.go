package v2

import (
	"context"
	"sync"

	"github.com/honestbank/kp/v2/internal/middleware"
)

type Processor[MessageType any] func(ctx context.Context, item *MessageType) error

type kp[MessageType any] struct {
	chain          middleware.Processor[*MessageType, error]
	shouldContinue bool
	mu             sync.Mutex
}

func (t *kp[MessageType]) getShouldContinue() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.shouldContinue
}

func (t *kp[MessageType]) AddMiddleware(middleware middleware.Middleware[*MessageType, error]) MessageProcessor[MessageType] {
	t.chain.AddMiddleware(middleware)

	return t
}

func (t *kp[MessageType]) Stop() {
	t.mu.Lock()
	t.shouldContinue = false
	t.mu.Unlock()
}

func (t *kp[MessageType]) Run(processor Processor[MessageType]) error {
	t.chain.AddMiddleware(middleware.FinalMiddleware[*MessageType, error](func(ctx context.Context, msg *MessageType) error {
		return processor(ctx, msg)
	}))

	for t.getShouldContinue() {
		ctx := context.Background()
		_ = t.chain.Process(ctx, nil)
	}

	return nil
}

func New[MessageType any]() MessageProcessor[MessageType] {
	return &kp[MessageType]{
		chain:          middleware.New[*MessageType, error](),
		shouldContinue: true,
	}
}
