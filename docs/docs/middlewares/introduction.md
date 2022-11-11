---
sidebar_position: 1
---

# Introduction
Middlewares are the core of KP, almost every feature except the retry and deadletters are built with middlewares.
They help us write feature in isolation, it makes it easier for testing and maintaining.

Middleware come in chain, and they're called in the order they were added.
Every middleware have ability to halt the operation completely.

## Configuration {#configuration}
To add middlewares to kp, simply call `AddMiddleware` and they'll be called in the order they were added.

```go
package main

import (
	"context"
	"fmt"
	"time"

	v2 "github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/middlewares"
)

type UserLoggedInEvent struct {
	UserID string
}

func main() {
	applicationName := "send-login-notification-worker"
	retryCount := 10
	kp := v2.New[UserLoggedInEvent]("user-logged-in", applicationName).
		WithRetryOrPanic("send-login-notification-retries", retryCount).
		WithDeadletterOrPanic("send-login-notification-failures")
	kp.AddMiddleware(middlewares.Measure("", ""))
	kp.AddMiddleware(middlewares.Tracing())
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
```
