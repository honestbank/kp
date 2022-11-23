---
sidebar_position: 2
---

# Retries
Retry feature sends failed messages into a retry topic and processes it when it receives the message.

We send the message to the worker's retry topic to avoid delays in retries causing the consumer to lag.

:::warning
Before you execute the `Run` method of KP, the retry topic must exist in Kafka.
Failing to do so will produce the message to the retry topic, but it won't be attempted till next start of KP.
:::

### Example {#example}

Take the worker and make a simple function call to add retries like the following:

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
	ensureTopicExists("send-login-notification-retries")
	kp := v2.New[UserLoggedInEvent]("user-logged-in", getConfig())
	kp.WithRetryOrPanic("send-login-notification-retries", 10) // + this line adds 10 retries
	err := kp.Process(processUserLoggedInEvent)
	if err != nil {
		panic(err) // do better error handling
	}
}

func getConfig() any {
    return nil // return config
}

func ensureTopicExists(topic string) {
    // check if the topic exists and create/error if it doesn't exist already
}

func processUserLoggedInEvent(ctx context.Context, message UserLoggedInEvent) error {
	// here, you can focus on your business logic.
	fmt.Printf("processing %v\n", message)
	time.Sleep(time.Millisecond * 200) // simulate long running process
	return nil // or error
}
```

When a message exceeds the configured retry limit, it ignores the message by default. See [next page](./deadletters.md) to configure deadletters.
