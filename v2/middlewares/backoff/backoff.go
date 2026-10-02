package backoff

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	backoff_policy "github.com/honestbank/backoff-policy"
	"github.com/honestbank/kp/v2/middlewares"
)

type backoff struct {
	p backoff_policy.BackoffPolicy
	// stopped is done after Stop, so that a wait ends at once.
	stopped context.Context
	stop    context.CancelFunc
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
		defer context.AfterFunc(b.stopped, cancel)()
		p.ExecuteWithContext(waitCtx, work)

		return err
	}
	b.p.Execute(work)

	return err
}

// Stop ends the wait in progress and every later wait, so that the message in progress runs at once.
// MessageProcessor.Stop calls it. It works when the policy is a backoff_policy.ContextBackoffPolicy, as the
// policies of backoff_policy.NewBackoff and backoff_policy.NewExponentialBackoffPolicy are.
func (b *backoff) Stop() {
	b.stop()
}

func NewBackoffMiddleware(policy backoff_policy.BackoffPolicy) middlewares.KPMiddleware[*kafka.Message] {
	stopped, stop := context.WithCancel(context.Background())

	return &backoff{p: policy, stopped: stopped, stop: stop}
}
