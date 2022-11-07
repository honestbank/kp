package middlewares

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"

	"github.com/honestbank/kp/v2/internal/middleware"
)

var operationDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "kp_operation_time_milliseconds",
	Buckets: []float64{1, 5, 50, 200, 500, 1_000, 2_000, 5_000, 15_000, 45_000},
}, []string{"result", "error"})

var totalProcessed = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "message_processed_total",
		Help: "Number of messages processed.",
	},
	[]string{"result", "error"},
)

type measurementMiddleware struct {
	pushClient *push.Pusher
}

func (m measurementMiddleware) SetupBackgroundJob() {
	go func() {
		err := m.pushClient.Push()
		if err != nil {
			fmt.Printf("error pushing metrics to gateway: %v\n", err)
		}
		time.Sleep(time.Second * 5)
	}()
}

func (m measurementMiddleware) Process(ctx context.Context, item *kafka.Message, next func(ctx context.Context, item *kafka.Message) error) error {
	startTime := time.Now()
	err := next(ctx, item)
	if err != nil {
		totalProcessed.WithLabelValues("failure", err.Error()).Inc()
		operationDuration.WithLabelValues("failure", err.Error()).Observe(float64(time.Since(startTime).Milliseconds()))

		return err
	}
	operationDuration.WithLabelValues("success", "").Observe(float64(time.Since(startTime).Milliseconds()))
	totalProcessed.WithLabelValues("success", "").Inc()

	return err
}

func Measure(gatewayURL string, applicationName string) middleware.Middleware[*kafka.Message, error] {
	pushClient := push.New(gatewayURL, applicationName).
		Grouping("framework", "kp").
		Collector(operationDuration).
		Collector(totalProcessed)

	mw := measurementMiddleware{
		pushClient: pushClient,
	}
	mw.SetupBackgroundJob()

	return mw
}
