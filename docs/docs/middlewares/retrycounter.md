---
sidebar_position: 5
---

# Retry Count
This is an informational middleware.
The purpose of this middleware is to provide an API to get the current retry count for a message.

:::info
If you forget to add the middleware and tried to access the value, it's going to return `int(0)`
:::

### Example {#example}

Simply add the `middlewares.RetryCount` before any other middlewares that depend on this information to be present.

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
	kp := v2.New[UserLoggedInEvent]("user-logged-in", getConfig())
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

func getConfig() any {
	return nil // return your config
}
```
