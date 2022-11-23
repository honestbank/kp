---
sidebar_position: 3
---

# Deadletters
If a particular message fails the configured number of times, it's simply ignored. But most of the time that's not what we want. In this example, it sends the message to a deadletter topic from where you can send alerts and check why they failed.

:::info
While you absolutely can use retries without deadletters, it'll probably be hard to setup re-processing of failed items. Using deadletters is highly recommended.
:::

### Example {#example}

Simply call the `WithDeadletter` or fluent `WithDeadletterOrPanic` method to enable deadletters

```go
package main

import (
	"context"
	"fmt"
	"time"

	v2 "github.com/honestbank/kp/v2"
)

type UserLoggedInEvent struct {
	UserID string
}

func main() {
	retryCount := 10
	kp := v2.New[UserLoggedInEvent]("user-logged-in", getConfig()).WithRetryOrPanic("send-login-notification-retries", retryCount)
	kp.WithDeadletterOrPanic("send-login-notification-failures") // + this line sends messages to deadletter topic after retry has exhausted it's limits
	err := kp.Process(processUserLoggedInEvent)
	if err != nil {
		panic(err) // do better error handling
	}
}

func getConfig() error {
    return nil // return config
}

func processUserLoggedInEvent(ctx context.Context, message UserLoggedInEvent) error {
	// here, you can focus on your business logic.
	fmt.Printf("processing %v\n", message)
	time.Sleep(time.Millisecond * 200) // simulate long running process
	return nil // or error
}
```

:::tip
Both deadletter and retries can be used independently or combined.
:::
