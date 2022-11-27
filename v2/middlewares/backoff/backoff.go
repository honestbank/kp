package backoff

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	backoff_policy "github.com/honestbank/backoff-policy"
	"github.com/honestbank/kp/v2/internal/middleware"
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

func NewBackoffMiddleware(policy backoff_policy.BackoffPolicy) middleware.Middleware[*kafka.Message, error] {
	return &backoff{p: policy}
}
