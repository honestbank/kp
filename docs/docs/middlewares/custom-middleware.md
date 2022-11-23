---
sidebar_position: 6
---

# Write your own
There could be cases that we've not covered in KP. In those scenarios you can write your own middleware. It's very easy.

All Middlewares implement the `Middleware[*kafka.Message, error]` interface, so that's all we have to do.

Definition of Middleware can be found [here](https://github.com/honestbank/kp/blob/52ed4f94b682835508513368314962f55d59fd1b/v2/internal/middleware/middleware.go#L19-L21)

:::tip
No need to worry about IN and OUT, simply write a `Middleware[*kafka.Message, error]` implementation.
:::

### Creating a log middleware {#implementation}
Let's say we needed to log as soon as retry count reaches a certain threshold to stdout

We'll need the following:
- threshold

Let's create a struct that can hold the above data and a function that accepts input for the values and finally implement the interface.

```go
package logmw

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/honestbank/kp/v2/internal/middleware"
	"github.com/honestbank/kp/v2/middlewares"
)

type logMw struct {
	Threshold int
}

func LogMiddleware(threshold int) middleware.Middleware[*kafka.Message, error] {
	return logMw{
		Threshold: threshold,
	}
}
func (r logMw) Process(ctx context.Context, item *kafka.Message, next func(ctx context.Context, item *kafka.Message) error) error {
	count := middlewares.RetryCountFromContext(ctx)
	r.logIfNeeded(count)
	return next(ctx, item)
}

func (r logMw) logIfNeeded(count int) {
	if count < r.Threshold {
		return
	}
	fmt.Printf("attempt %d reached\n", count)
}
```


### Usage {#usage}

:::warning
Because the `RetryCountFromContext` is being used, `middlewares.RetryCount` needs to be added before the `LogMiddleware`.
:::

```go
package main

import (
	"context"

	"github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/middlewares"
)

type UserLoggedInEvent struct {
	UserID string
}

func main() {
	defer setupTracing()() // this is important and not included in kp by default
	kp := v2.New[UserLoggedInEvent]("user-logged-in", getConfig())
	kp.WithRetryOrPanic("send-login-notification-retries", 10)
	kp.AddMiddleware(middlewares.RetryCount())
	kp.AddMiddleware(logmw.LogMiddleware())
	err := kp.Process(processUserLoggedInEvent)
	if err != nil {
		panic(err) // do better error handling
	}
}

func processUserLoggedInEvent(ctx context.Context, message UserLoggedInEvent) error {
	// here, you can focus on your business logic.
	fmt.Printf("processing %v\n", message)
	time.Sleep(time.Millisecond * 200) // simulate long running process
	return nil                         // or error
}

func getConfig() any {
	return nil // return your config
}
```
