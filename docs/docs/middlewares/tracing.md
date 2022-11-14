---
sidebar_position: 4
---

# Tracing
Chances are if you have kafka, you also want a distributed tracing. Because tracing is so common, we've included tracing middleware.

We don't control where the traces go to or in which format, we simply create spans.

:::info
If you use tracing middleware, you need to set up a trace provider yourself. That way you can send traces to any endpoint on any format. We produce spans using [opentelemetry](https://pkg.go.dev/go.opentelemetry.io/otel)
:::

## Sample trace {#screenshot}
In the following a message failed to be processed 6 times and it successfully processed the message the 7th time.
![tracing screenshot](../../static/img/tracing_example.png)

### Configuration {#configuration}

First, you'll need a trace provider only then start the processor with the tracing middleware

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
	defer setupTracing()() // this is important and not included in kp
	applicationName := "send-login-notification-worker"
	kp := v2.New[UserLoggedInEvent]("user-logged-in", applicationName)
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
```
