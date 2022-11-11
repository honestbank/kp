---
sidebar_position: 3
---
# Measurements
Every service needs some kind of measurements. KP comes with a prometheus push metric middleware.

:::warning
Measurement middleware should come after backoff middleware so that we don't measure wait times (but you can).
:::

### Configuration {#configuration}

Take the processor and make a simple function call to add retries like the following:

```go
package main

import (
	"context"
	"fmt"
	"time"

	backoff_policy "github.com/honestbank/backoff-policy"
	"github.com/honestbank/backoff-policy/policies"
	v2 "github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/middlewares"
)

type UserLoggedInEvent struct {
	UserID string
}

func main() {
	applicationName := "send-login-notification-worker"
	kp := v2.New[UserLoggedInEvent]("user-logged-in", applicationName)
	kp.WithRetryOrPanic("send-login-notification-retries", 10)
	kp.WithDeadletterOrPanic("send-login-notification-failures")
	kp.AddMiddleware(middlewares.Measure("http://path/to/push/gateway", applicationName)) // simply add a measurement middleware to get free metrics
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
```
