package retry_count

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/honestbank/kp/v2/internal/middleware"
	"github.com/honestbank/kp/v2/internal/retrycounter"
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

func NewRetryCountMiddleware() middleware.Middleware[*kafka.Message, error] {
	return retryCountMw{}
}
