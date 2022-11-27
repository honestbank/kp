---
sidebar_position: 1
---

# Introduction
Middlewares are the core of KP, almost every feature except the retry and deadletters are built with middlewares.

They help us write isolated, testable and maintainable features.

Middleware come in chain, and they're called in the order they were added.
Every middleware have ability to halt the operation completely.

## Example {#example}
To add middleware to KP, simply call `AddMiddleware` and they'll be called in the order they were added.

```go
package main

import (
	"context"
	"fmt"
	"time"

	v2 "github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/middlewares/measurement"
	"github.com/honestbank/kp/v2/middlewares/tracing"
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
	kp.AddMiddleware(tracing.NewTracingMiddleware())
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
