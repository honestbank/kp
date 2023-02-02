package backoff

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	backoff_policy "github.com/honestbank/backoff-policy"
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
