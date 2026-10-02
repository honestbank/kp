package backoff

import (
	"context"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	backoff_policy "github.com/honestbank/backoff-policy"
	"github.com/honestbank/backoff-policy/policies"
	"github.com/honestbank/kp/v2/middlewares"
)

type backoff struct {
	p backoff_policy.BackoffPolicy
}

func (b backoff) Process(ctx context.Context, item *kafka.Message, next func(ctx context.Context, item *kafka.Message) error) error {
	var err error
	b.p.Execute(func(marker backoff_policy.Marker) {
		err = next(ctx, item)
		if err != nil {
			marker.MarkFailure()
			return
		}
		marker.MarkSuccess()
	})

	return err
}

func NewBackoffMiddleware(policy backoff_policy.BackoffPolicy) middlewares.KPMiddleware[*kafka.Message] {
	return &backoff{p: policy}
}

// InterruptibleBackoff waits as the backoff middleware does, but Stop or a done context ends the wait. The backoff
// middleware sleeps in backoff_policy.Execute, and nothing can end that sleep, so a shutdown has to wait for it.
type InterruptibleBackoff struct {
	policy   policies.Policy
	failures int
	stop     chan struct{}
	stopOnce sync.Once
}

// NewInterruptibleBackoffMiddleware returns a backoff middleware that waits policy(failures) before each message. A
// failed message adds 1 to the failures, and a successful message takes 1 away, as backoff_policy.NewBackoff does.
// Call Stop next to MessageProcessor.Stop, so that the message in progress runs at once.
func NewInterruptibleBackoffMiddleware(policy policies.Policy) *InterruptibleBackoff {
	return &InterruptibleBackoff{policy: policy, stop: make(chan struct{})}
}

// Process waits for the backoff of the failures so far, then runs next. kp processes 1 message at a time, so the
// failures need no lock.
func (b *InterruptibleBackoff) Process(ctx context.Context, item *kafka.Message, next func(ctx context.Context, item *kafka.Message) error) error {
	if wait := b.policy(b.failures); wait > 0 {
		timer := time.NewTimer(wait)
		select {
		case <-timer.C:
		case <-b.stop:
			timer.Stop()
		case <-ctx.Done():
			timer.Stop()
		}
	}

	err := next(ctx, item)
	if err != nil {
		b.failures++
	} else if b.failures > 0 {
		b.failures--
	}

	return err
}

// Stop ends the wait in progress and every later wait. It is safe to call more than once.
func (b *InterruptibleBackoff) Stop() {
	b.stopOnce.Do(func() {
		close(b.stop)
	})
}
