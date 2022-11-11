---
sidebar_position: 1
---

# Why use KP?
Kafka is a message bus where a message is either committed or not, without committing a message, we can't process the next message.

Implementing retries, backoffs and tracing becomes a little difficult as they all pollute your business code.

KP provides everything out of the box. With one line each, you can configure backoffs, retries, deadletters, tracing, measurements etc. It makes it in a way that your business logic is a simple user-defined-function free of all the other logic which isn't the core of the business.

Take the following example:

```go
package main

import (
	"context"
	"time"

	backoff_policy "github.com/honestbank/backoff-policy"
	v2 "github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/middlewares"
)

type UserLoggedInEvent struct {
	UserID string
}

func main() {
	initializeTracer()
	retryCount := 10
	applicationName := "send-login-notification-worker"
	kp := v2.New[UserLoggedInEvent]("user-logged-in", applicationName)
	kp.
		WithRetryOrPanic("send-login-notification-retries", retryCount). // 1 line to enable retries
		WithDeadletterOrPanic("send-login-notification-failures"). // 1 line to enable deadlettering
		AddMiddleware(middlewares.Backoff(backoff_policy.NewExponentialBackoffPolicy(time.Millisecond*200, 10))). // 1 line to enable backoffs
		AddMiddleware(middlewares.Tracing()). // 1 line to enable tracing
		AddMiddleware(middlewares.Measure("path/to/prometheus-push-gateway", applicationName)). // 1 line to enable measurements
		AddMiddleware(middlewares.RetryCount())
	kp.Process(func(ctx context.Context, message UserLoggedInEvent) error {
		// process the message and return error if it fails. don't worry about retries here.
		// but if you need to know the current count use the following:
		count := middlewares.RetryCountFromContext(ctx)
		print(count)
		return nil
    })
}

func initializeTracer() {
	// initialize your tracer
}
```
In the above example, we added retries, deadletters, backoffs, tracing and measurements all taking 1 line each.

:::tip
It is also possible to write your own middleware if kp doesn't cover usecase. And it's very simple to do so. Checkout writing middlewares page to learn how.
:::