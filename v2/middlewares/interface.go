package middlewares

import "github.com/honestbank/kp/v2/internal/middleware"

type KPMiddleware[T any] interface {
	middleware.Middleware[T, error]
}
