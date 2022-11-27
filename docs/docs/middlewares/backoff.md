---
sidebar_position: 2
---
# Backoff
When the worker starts receiving failures, ideally, it should slow down the message processing and let the underlying services recover. Doing this is as simple as adding 1 middleware like in the following example:

:::tip
You don't have to use the exponential backoff as shown in the example, you can bring your own policy as well.
:::

### Example {#example}

:::warning
Adding multiple backoff middleware is possible.
:::

```go
package main

import (
	"context"
	"fmt"
	"time"

	backoff_policy "github.com/honestbank/backoff-policy"
	"github.com/honestbank/backoff-policy/policies"
	v2 "github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/middlewares/backoff"
)

type UserLoggedInEvent struct {
	UserID string
}

func main() {
	kp := v2.New[UserLoggedInEvent]("user-logged-in", getConfig())
	kp.WithRetryOrPanic("send-login-notification-retries", 10) // + this line adds 10 retries
	exponent, duration, maxBackoffCount := 1.5, time.Millisecond*200, 10
	backoffPolicy := backoff_policy.NewBackoff(policies.GetExponentialPolicy(exponent, duration, maxBackoffCount))
	kp.AddMiddleware(backoff.NewBackoffMiddleware(backoffPolicy)) // simply add a backoff middleware to back off.
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
