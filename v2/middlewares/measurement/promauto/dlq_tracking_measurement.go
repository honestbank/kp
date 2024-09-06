package promauto

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/honestbank/kp/v2/internal/retrycounter"
	"github.com/honestbank/kp/v2/middlewares"
)

var retryCountReachedDuration = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "kp_retry_count_reached_total",
}, []string{"application_name", "topic"})

type dlqTrackingMeasurementMiddleware struct {
	applicationName string
	topic           string
	threshold       int
}

func (m dlqTrackingMeasurementMiddleware) Process(ctx context.Context, item *kafka.Message, next func(ctx context.Context, item *kafka.Message) error) error {
	if retrycounter.GetCount(item) >= m.threshold {
		retryCountReachedDuration.WithLabelValues(m.applicationName, m.topic).Inc()
	}
	return next(ctx, item)
}

func NewDLQTrackingMeasurementMiddleware(applicationName, topic string, threshold int) middlewares.KPMiddleware[*kafka.Message] {
	return dlqTrackingMeasurementMiddleware{
		applicationName: applicationName,
		topic:           topic,
		threshold:       threshold,
	}
}
