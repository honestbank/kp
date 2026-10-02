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

### Interruptible backoff {#interruptible-backoff}

The backoff middleware sleeps, and nothing can end that sleep. After a few failures, a shutdown must wait for the full backoff before the consumer stops.
`NewInterruptibleBackoffMiddleware` waits the same way, but `Stop` (or a done context) ends the wait, and the message in progress then runs at once.
It takes a `policies.Policy` and counts the failures itself, as `backoff_policy.NewBackoff` does.

```go
processor := v2.New[kafka.Message]()
backoffMiddleware := backoff.NewInterruptibleBackoffMiddleware(policies.GetExponentialPolicy(1.5, time.Millisecond*200, 10))
processor.AddMiddleware(backoffMiddleware)

stop := func() {
	processor.Stop()
	backoffMiddleware.Stop() // the message in progress does not wait for the backoff
}
```
