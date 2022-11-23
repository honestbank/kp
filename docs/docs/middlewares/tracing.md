---
sidebar_position: 4
---

# Tracing
Chances are if you're using Kafka, you'll need distributed tracing. KP comes with a tracing middleware.

Tracing middleware does not control the format or the destination, it simply creates spans. Setting the destination and format is done by setting a trace provider.

:::warning
Using trace middleware without a correctly configured trace provider will result in invalid spans being produced.
:::

## Sample trace {#screenshot}
In the following picture, a message failed to be processed 6 times and was successfully processed the 7th time.
![tracing screenshot](../../static/img/tracing_example.png)

### Example {#example}

Configure trace provider and add `middlewares.Tracing` to enable traces.

```go
package main

import (
	"context"
	"os"

	"github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/middlewares"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

type UserLoggedInEvent struct {
	UserID string
}

func main() {
	defer setupTracing()() // this is important and not included in KP
	kp := v2.New[UserLoggedInEvent]("user-logged-in", getConfig())
	kp.WithRetryOrPanic("send-login-notification-retries", 10)
	kp.AddMiddleware(middlewares.Tracing()) // This adds tracing middleware
	err := kp.Process(processUserLoggedInEvent)
	if err != nil {
		panic(err) // do better error handling
	}
}

func processUserLoggedInEvent(ctx context.Context, message UserLoggedInEvent) error {
	// here, you can focus on your business logic.
	fmt.Printf("processing %v\n", message)
	time.Sleep(time.Millisecond * 200) // simulate long running process
	return nil                         // or error
}

func setupTracing() func() {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(os.Getenv("COLLECTOR_URL"))))
	if err != nil {
		panic(err)
	}
	resc := resource.Default()
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resc),
		trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("form-service-web"))),
	)
	otel.SetTracerProvider(tp)
	return func() {
		exporter.Shutdown(context.Background())
		tp.Shutdown(context.Background())
	}
}

func getConfig() any {
	return nil // return your config
}
```
