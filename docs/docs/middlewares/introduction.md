---
sidebar_position: 1
---

# Introduction
Middleware is a software pattern that allows you to intercept and modify the processing of messages in a chain. It is a powerful tool for customizing the behavior of a system and adding additional functionality, such as logging, retries, deadlettering, backoffs, and tracing.

In the KP library, every feature is implemented as a middleware. This means that you can add or remove features by adding or removing middlewares from the middleware chain.

## Adding Middleware to the KP Library
The KP library provides a way to add middlewares to the middleware chain through the `AddMiddleware`

This method takes a middleware as an argument and returns a new MessageProcessor with the middleware added to the chain. Here is an example of how to use the `AddMiddleware` method:

## Example {#example}

```go
package main

import (
	"context"
	"fmt"
	"time"

	v2 "github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/middlewares/measurement"
	"github.com/honestbank/kp/v2/middlewares/tracing"
	"go.opentelemetry.io/otel"
)

type UserLoggedInEvent struct {
	UserID string
}

func main() {
	retryCount := 10
	kp := v2.New[UserLoggedInEvent]("user-logged-in", getConfig()).
		WithRetryOrPanic("send-login-notification-retries", retryCount).
		WithDeadletterOrPanic("send-login-notification-failures")
	kp.AddMiddleware(measurement.NewMeasurementMiddleware("/path/to/push-gateway", "application-name"))
	tp := otel.GetTracerProvider() // you'll need to set the tracer provider as well
	kp.AddMiddleware(tracing.NewTracingMiddleware(tp))
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

func getConfig() any {
	return nil // return your config
}
```
### Notes {#notes}

- Middlewares are executed in the order they are added to the chain. It is important to consider the order in which middlewares are added.
- The KP library provides several built-in middlewares for common features, such as consuming messages, retries, deadletters, backoffs, and tracing. You can also write custom middlewares to meet your specific

