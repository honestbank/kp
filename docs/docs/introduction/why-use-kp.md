---
sidebar_position: 1
---

# Why use KP?
Kafka is a message bus where a message is either committed or not, without committing a message, we can't process the next message.

Implementing retries, backoffs and tracing becomes a little difficult as they all pollute your business code.

KP provides everything out of the box. With one line each, you can configure backoffs, retries, deadletters, tracing, measurements etc. It makes it in a way that your business logic is a simple user-defined-function free of all the other logic which isn't the core of the business.

Take the following example:

:::tip
Please check [this page](../introduction/configuration.md) for detailed configuration option
:::

```go
package main

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/honestbank/kp/v2/middlewares/consumer"
	"github.com/honestbank/kp/v2/middlewares/deadletter"
	"github.com/honestbank/kp/v2/middlewares/retry"
	"time"

	backoff_policy "github.com/honestbank/backoff-policy"
	v2 "github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/middlewares/backoff"
	"github.com/honestbank/kp/v2/middlewares/measurement"
	"github.com/honestbank/kp/v2/middlewares/retry_count"
	"github.com/honestbank/kp/v2/middlewares/tracing"
	"go.opentelemetry.io/otel"
)

type UserLoggedInEvent struct {
	UserID string
}

func main() {
	initializeTracer()
	retryCount := 10
	applicationName := "send-login-notification-worker"
	kp := v2.New[kafka.Message]()
	kp.
		AddMiddleware(consumer.NewConsumerMiddleware(getConsumer())).
		AddMiddleware(deadletter.NewDeadletterMiddleware(getDeadletterProducer(), 10, func(err error) {})).
		AddMiddleware(retry.NewRetryMiddleware(getRetryMiddleware(), func(err error) {})).
		AddMiddleware(backoff.NewBackoffMiddleware(backoff_policy.NewExponentialBackoffPolicy(time.Millisecond*200, 10))). // 1 line to enable backoffs
		AddMiddleware(tracing.NewTracingMiddleware(otel.GetTracerProvider())). // 1 line to enable tracing
		AddMiddleware(measurement.NewMeasurementMiddleware("path/to/prometheus-push-gateway", applicationName)). // 1 line to enable measurements
		AddMiddleware(retry_count.NewRetryCountMiddleware())
	kp.Process(func(ctx context.Context, message *kafka.Message) error {
		// process the message and return error if it fails. don't worry about retries here.
		// but if you need to know the current count use the following:
		count := retry_count.FromContext(ctx)
		print(count)
		return nil
	})
}

func getConfig() any {
	return nil // return your config
}

func initializeTracer() {
	// initialize your tracer
}
```
In the above example, we added retries, deadletters, backoffs, tracing and measurements all taking 1 line each.

:::tip
It is also possible to write your own middleware if KP doesn't cover usecase. And it's very simple to do so. Checkout writing middlewares page to learn how.
:::
