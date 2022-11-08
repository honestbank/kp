package middlewares

import (
	"context"
	"errors"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.opentelemetry.io/otel"

	"github.com/honestbank/kp/v2/internal/kafkaheaders"
	"github.com/honestbank/kp/v2/internal/middleware"
	"github.com/honestbank/kp/v2/internal/tracing"
)

type tracingMw struct{}

func (t tracingMw) Process(ctx context.Context, item *kafka.Message, next func(ctx context.Context, item *kafka.Message) error) error {
	tracer := otel.GetTracerProvider().Tracer("empty")
	ctx, span := tracer.Start(tracing.ExtractTraceContext(ctx, item), "process")
	res := next(ctx, item)
	span.End()
	// we want to close before re-producing the message
	if kafkaheaders.Get("traceparent", item) == nil { // let's ask.
		tracing.InjectTraceHeaders(ctx, item)
	}

	return res
}

func Tracing() (middleware.Middleware[*kafka.Message, error], error) {
	tp := otel.GetTracerProvider()
	if tp == nil {
		return nil, errors.New("trace provider is not set")
	}
	return tracingMw{}, nil
}
