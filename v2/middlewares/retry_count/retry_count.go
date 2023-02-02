package retry_count

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/honestbank/kp/v2/internal/retrycounter"
	"github.com/honestbank/kp/v2/middlewares"
)

type ctxKey struct{}

func FromContext(ctx context.Context) int {
	val := ctx.Value(&ctxKey{})
	if val == nil {
		return 0
	}

	return val.(int) // it's not possible for this val not to be int since ctxKey is unexported
}

type retryCountMw struct {
}

func (r retryCountMw) Process(ctx context.Context, item *kafka.Message, next func(ctx context.Context, item *kafka.Message) error) error {
	retryCount := retrycounter.GetCount(item)
	return next(context.WithValue(ctx, &ctxKey{}, retryCount), item)
}

func NewRetryCountMiddleware() middlewares.KPMiddleware[*kafka.Message] {
	return retryCountMw{}
}
