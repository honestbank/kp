---
sidebar_position: 4
---
# Backoff
The backoff middleware slows down the processing of messages when errors occur. It uses a backoff policy to determine how long to wait between processing attempts.

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

	"github.com/confluentinc/confluent-kafka-go/kafka"
	backoff_policy "github.com/honestbank/backoff-policy"
	"github.com/honestbank/backoff-policy/policies"
	v2 "github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/middlewares/backoff"
)

func main() {
	kp := v2.New[kafka.Message]()
	exponent, duration, maxBackoffCount := 1.5, time.Millisecond*200, 10
	backoffPolicy := backoff_policy.NewBackoff(policies.GetExponentialPolicy(exponent, duration, maxBackoffCount))
	kp.AddMiddleware(backoff.NewBackoffMiddleware(backoffPolicy)) // simply add a backoff middleware to back off.
	err := kp.Process(processUserLoggedInEvent)
	if err != nil {
		panic(err) // do better error handling
	}
}

func processUserLoggedInEvent(ctx context.Context, message *kafka.Message) error {
	// here, you can focus on your business logic.
	fmt.Printf("processing %v\n", message)
	time.Sleep(time.Millisecond * 200) // simulate long running process
	return nil // or error
}

func getConfig() any {
	return nil // return your config
}
```

### Stop ends the wait {#stop-ends-the-wait}

`MessageProcessor.Stop` ends the backoff wait in progress, and every later wait, so the message in progress runs at once and `Run` returns without the wait.
This works when the policy is a `backoff_policy.ContextBackoffPolicy`, as the policies of `backoff_policy.NewBackoff` and `backoff_policy.NewExponentialBackoffPolicy` are (backoff-policy v1.4.0 or later).
A done message context also ends the wait. A policy with `Execute` only waits as before.

```go
processor := v2.New[kafka.Message]()
processor.AddMiddleware(backoff.NewBackoffMiddleware(backoff_policy.NewExponentialBackoffPolicy(time.Millisecond*200, 10)))

// on shutdown: the message in progress does not wait for the backoff
processor.Stop()
```
