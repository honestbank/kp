---
sidebar_position: 1
---

# Retries
When there's an error while processing a message, most of the time we want to try again.
But we might not want to try at the same time causing the consumer to lag.
Retry feature sends failed messages into a retry topic and processes it when it receives the message.

:::warning
Before you start the `Run` method of kp, retry topic must exist in kafka.
Failing to do so will produce the message to the retry topic, but it won't be attempted till next start of kp.
:::

### Configuration {#configuration}

Take the processor and make a simple function call to add retries like the following:

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
	applicationName := "send-login-notification-worker"
	kp := v2.New[UserLoggedInEvent]("user-logged-in", applicationName)
	kp.WithRetryOrPanic("send-login-notification-retries", 10) // + this line adds 10 retries
	err := kp.Process(processUserLoggedInEvent)
	if err != nil {
		panic(err) // do better error handling
	}
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
