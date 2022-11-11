---
sidebar_position: 5
---

# Retry Count
This is an informational middleware. It doesn't change the execution in any way.
The purpose of this middleware is to provide an API to get the retry count from kafka message.
If we're seeing the message for the first time, it sets the value as 0.

:::info
If you forget to add the middleware and tried to access the value, it's going to return `int(0)`
:::

### Configuration {#configuration}

First, you'll need an instance with retry enabled. Then simply add the middleware anywhere in the chain.

```go
package main

import (
	"context"
	"fmt"

	"github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/middlewares"
)

type UserLoggedInEvent struct {
	UserID string
}

func main() {
	applicationName := "send-login-notification-worker"
	kp := v2.New[UserLoggedInEvent]("user-logged-in", applicationName)
	kp.WithRetryOrPanic("send-login-notification-retries", 10)
	kp.AddMiddleware(middlewares.RetryCount()) // This adds retry count middleware
	err := kp.Process(processUserLoggedInEvent)
	if err != nil {
		panic(err) // do better error handling
	}
}

func processUserLoggedInEvent(ctx context.Context, message UserLoggedInEvent) error {
	count := middlewares.RetryCountFromContext(ctx) // this is how to get the count
	fmt.Printf("this message for %s user was processed %d times", message.UserID, count)
	return nil // or error
}
```
