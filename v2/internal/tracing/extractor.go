package tracing

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.opentelemetry.io/otel/propagation"
)

func ExtractTraceContext(ctx context.Context, message *kafka.Message) context.Context {
	mapHeaders := make(map[string]string)
	for _, v := range message.Headers {
		mapHeaders[v.Key] = string(v.Value)
	}
	traceCTX := propagation.TraceContext{}
	ctx = traceCTX.Extract(ctx, propagation.MapCarrier(mapHeaders))

	return ctx
}
