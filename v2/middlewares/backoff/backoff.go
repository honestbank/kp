package backoff

import (
	"context"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	backoff_policy "github.com/honestbank/backoff-policy"
	"github.com/honestbank/kp/v2/middlewares"
)

type backoff struct {
	p backoff_policy.BackoffPolicy
	// mu guards stopped and cancelWait, because Stop runs in another goroutine.
	mu      sync.Mutex
	stopped bool
	// cancelWait ends the wait in progress. It is nil when no message waits.
	cancelWait context.CancelFunc
}

func (b *backoff) Process(ctx context.Context, item *kafka.Message, next func(ctx context.Context, item *kafka.Message) error) error {
	var err error
	work := func(marker backoff_policy.Marker) {
		err = next(ctx, item)
		if err != nil {
			marker.MarkFailure()

			return
		}
		marker.MarkSuccess()
	}

	// A policy that can end its wait ends it when ctx is done or the middleware stops. next still gets ctx.
	if p, ok := b.p.(backoff_policy.ContextBackoffPolicy); ok {
		waitCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		b.startWait(cancel)
		defer b.endWait()
		p.ExecuteWithContext(waitCtx, work)

		return err
	}
	b.p.Execute(work)

	return err
}

// startWait keeps cancel for Stop. After Stop, it cancels at once, so the message does not wait.
func (b *backoff) startWait(cancel context.CancelFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.stopped {
		cancel()

		return
	}
	b.cancelWait = cancel
}

func (b *backoff) endWait() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.cancelWait = nil
}

// Stop ends the wait in progress and every later wait, so that the message in progress runs at once.
// MessageProcessor.Stop calls it. It works when the policy is a backoff_policy.ContextBackoffPolicy, as the
// policies of backoff_policy.NewBackoff and backoff_policy.NewExponentialBackoffPolicy are.
func (b *backoff) Stop() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.stopped = true
	if b.cancelWait != nil {
		b.cancelWait()
	}
}

func NewBackoffMiddleware(policy backoff_policy.BackoffPolicy) middlewares.KPMiddleware[*kafka.Message] {
	return &backoff{p: policy}
}
