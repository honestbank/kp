package promauto

import (
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/honestbank/kp/v2/middlewares"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var operationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "kp_operation_duration_milliseconds",
	Buckets: []float64{1, 5, 50, 200, 500, 1_000, 2_000, 5_000, 15_000, 45_000},
}, []string{"application_name", "result", "error"})

type measurementMiddleware struct {
	applicationName string
}

func (m measurementMiddleware) Process(ctx context.Context, item *kafka.Message, next func(ctx context.Context, item *kafka.Message) error) error {
	startTime := time.Now()
	err := next(ctx, item)
	if err != nil {
		operationDuration.WithLabelValues(m.applicationName, "failure", err.Error()).Observe(float64(time.Since(startTime).Milliseconds()))

		return err
	}
	operationDuration.WithLabelValues(m.applicationName, "success", "").Observe(float64(time.Since(startTime).Milliseconds()))

	return err
}

func NewMeasurementMiddleware(applicationName string) middlewares.KPMiddleware[*kafka.Message] {
	return measurementMiddleware{
		applicationName: applicationName,
	}
}
