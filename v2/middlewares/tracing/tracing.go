package tracing

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/honestbank/kp/v2/internal/tracing"
	"github.com/honestbank/kp/v2/kafkaheaders"
	"github.com/honestbank/kp/v2/middlewares"
)

type tracingMiddleware struct {
	tracerProvider *trace.TracerProvider
}

func (t *tracingMiddleware) Process(ctx context.Context, item *kafka.Message, next func(ctx context.Context, item *kafka.Message) error) error {
	tracer := t.tracerProvider.Tracer("kp")
	ctx, span := tracer.Start(tracing.ExtractTraceContext(ctx, item), "process")
	err := next(ctx, item)
	if err != nil {
		span.RecordError(err)
		// https://github.com/open-telemetry/opentelemetry-go/blob/2694dbfdba21d4f5496e873d7f69f04ada040af1/trace/trace.go#L356-L360
		// we need to set status our selves
		span.SetStatus(codes.Error, "error span")
	}
	span.End()
	// we want to close before re-producing the message
	if kafkaheaders.Get("traceparent", item) == nil { // let's ask.
		tracing.InjectTraceHeaders(ctx, item)
	}

	return err
}

func NewTracingMiddleware(tracerProvider *trace.TracerProvider) middlewares.KPMiddleware[*kafka.Message] {
	return &tracingMiddleware{
		tracerProvider: tracerProvider,
	}
}
