---
sidebar_position: 3
---
# Metrics
Every service should be able to notify of its performance/stability. KP comes with a prometheus push metrics middleware.

:::warning
Measurement middleware should come after the backoff middleware so that we don't measure wait times (but you can).
:::

### Example {#example}

Simply add `middlewares.Measure` to enable metrics tracking.

```go
package main

import (
	"context"
	"fmt"
	"time"

	backoff_policy "github.com/honestbank/backoff-policy"
	"github.com/honestbank/backoff-policy/policies"
	v2 "github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/middlewares/measurement"
)

type UserLoggedInEvent struct {
	UserID string
}

func main() {
	applicationName := "send-login-notification-worker"
	kp := v2.New[UserLoggedInEvent]("user-logged-in", getConfig())
	kp.WithRetryOrPanic("send-login-notification-retries", 10)
	kp.WithDeadletterOrPanic("send-login-notification-failures")
	kp.AddMiddleware(measurement.NewMeasurementMiddleware("http://path/to/push/gateway", applicationName)) // simply add a measurement middleware to get free metrics
	err := kp.Process(processUserLoggedInEvent)
	if err != nil {
		panic(err) // do better error handling
	}
}

func processUserLoggedInEvent(ctx context.Context, message UserLoggedInEvent) error {
	// here, you can focus on your business logic.
	fmt.Printf("processing %v\n", message)
	time.Sleep(time.Millisecond * 200) // simulate long running process
	return nil // or error
}

func getConfig() any {
	return nil // return your config
}
```
