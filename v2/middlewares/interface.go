package middlewares

import "context"

type KPMiddleware[T any] interface {
	Process(ctx context.Context, item T, next func(ctx context.Context, item T) error) error
}
