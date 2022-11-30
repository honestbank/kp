package tracing

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.opentelemetry.io/otel/propagation"

	"github.com/honestbank/kp/v2/kafkaheaders"
)

func InjectTraceHeaders(context context.Context, message *kafka.Message) {
	propagationHeadersMap := make(map[string]string)
	dummyCtx := propagation.TraceContext{}
	dummyCtx.Inject(context, propagation.MapCarrier(propagationHeadersMap))
	for key, value := range propagationHeadersMap {
		kafkaheaders.Set(message, key, value)
	}
}
