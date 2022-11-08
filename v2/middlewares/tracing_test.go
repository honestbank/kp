package middlewares_test

import (
	"context"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/honestbank/kp/v2/internal/kafkaheaders"
	"github.com/honestbank/kp/v2/middlewares"
)

func TestTracingMw_Process(t *testing.T) {
	t.Run("adds trace related header when not present", func(t *testing.T) {
		//// need to initialize trace provider
		tp := trace.NewTracerProvider()
		otel.SetTracerProvider(tp)
		message := &kafka.Message{Headers: []kafka.Header{}}
		tracing, err := middlewares.Tracing()
		assert.NoError(t, err)
		tracing.Process(context.Background(), message, func(ctx context.Context, message2 *kafka.Message) error {
			return nil
		})
		assert.Greater(t, len(message.Headers), 0)
	})

	t.Run("if a message already has traceID, it should result in same traceID", func(t *testing.T) {
		message := &kafka.Message{Headers: []kafka.Header{}}
		kafkaheaders.Set(message, "traceparent", "00-e191a9feec1f18ba0c0d82eb0830a7d8-c611513a9ed84e4d-01")
		tracing, err := middlewares.Tracing()
		assert.NoError(t, err)
		tracing.Process(context.Background(), message, func(ctx context.Context, message2 *kafka.Message) error {
			traceParent := *kafkaheaders.Get("traceparent", message2)
			assert.Equal(t, "00-e191a9feec1f18ba0c0d82eb0830a7d8-c611513a9ed84e4d-01", traceParent)
			return nil
		})
		assert.Equal(t, "00-e191a9feec1f18ba0c0d82eb0830a7d8-c611513a9ed84e4d-01", *kafkaheaders.Get("traceparent", message))
	})
}
