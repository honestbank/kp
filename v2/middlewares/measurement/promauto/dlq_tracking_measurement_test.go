package promauto_test

import (
	"context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/honestbank/kp/v2/internal/retrycounter"
	"github.com/honestbank/kp/v2/middlewares/measurement/promauto"
)

func TestDLQTrackingMeasurementMiddleware(t *testing.T) {
	t.Run("should not report metric if threshold not reached", func(t *testing.T) {
		message := &kafka.Message{}
		retrycounter.SetCount(message, 1)
		_ = promauto.NewDLQTrackingMeasurementMiddleware("mock-app", "mock-topic", 3).Process(context.Background(), message, func(ctx context.Context, item *kafka.Message) error {
			return nil
		})
		metrics, err := scrapeMetrics()
		require.NoError(t, err)
		assert.NotContains(t, metrics, `kp_retry_count_reached_total`)
	})
	t.Run("should not report metric if threshold not reached", func(t *testing.T) {
		message := &kafka.Message{}
		retrycounter.SetCount(message, 3)
		_ = promauto.NewDLQTrackingMeasurementMiddleware("mock-app", "mock-topic", 3).Process(context.Background(), message, func(ctx context.Context, item *kafka.Message) error {
			return nil
		})
		metrics, err := scrapeMetrics()
		require.NoError(t, err)
		assert.Contains(t, metrics, `kp_retry_count_reached_total{application_name="mock-app",topic="mock-topic"} 1`)
	})
}

func scrapeMetrics() (string, error) {
	recorder := httptest.NewRecorder()
	promhttp.Handler().ServeHTTP(recorder, httptest.NewRequest("GET", "/metrics", nil))
	body, err := io.ReadAll(recorder.Result().Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
