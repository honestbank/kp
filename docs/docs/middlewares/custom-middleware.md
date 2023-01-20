---
sidebar_position: 10
---

# Write your own
Middlewares are a powerful way to customize the behavior of the library by intercepting and modifying the processing of messages. They are implemented as implementations of the KPMiddleware interface, which has a single method:
```go
type KPMiddleware[T any] interface {
	Process(ctx context.Context, item T, next func(ctx context.Context, item T) error) error
}
```

Definition of Middleware can be found [here](https://github.com/honestbank/kp/blob/e5d8f1e0ec1cceaa898e6a1ed7f7b468021b7b74/v2/middlewares/interface.go#L5-L7)


The Process method is called for each message and receives the following arguments:
- `ctx`: a context object that can be used to cancel the processing of messages or pass values between middlewares.
- `item`: the message being processed.
- `next`: a function that represents the next middleware in the chain. It should be called to pass the message to the next middleware, or to terminate the chain if it is the last middleware.


### Example: Log Middleware {#implementation}
Here is an example of a log middleware that logs "message processed at" and the current time whenever a message is processed:

```go
package logmw

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/honestbank/kp/v2/internal/middleware"
	"github.com/honestbank/kp/v2/middlewares"
)

type logMiddleware struct{}

func (l logMiddleware) Process(ctx context.Context, item T, next func(ctx context.Context, item T) error) error {
	log.Println("message processed at", time.Now())
	return next(ctx, item)
}

func NewLogMiddleware() KPMiddleware[T] {
	return &logMiddleware{}
}
```


### Usage {#usage}
To use the log middleware simply add it like the following:

```go
package main

import (
	"context"
	consumer2 "github.com/honestbank/kp/v2/consumer"
	"github.com/honestbank/kp/v2/middlewares/consumer"

	"github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/middlewares/retry_count"
)

func getConsumer() consumer2.Consumer {
    return nil
}

func main() {
	defer setupTracing()() // this is important and not included in kp by default
	kp := v2.New[*kafka.Message]()
	kp.AddMiddleware(consumer.NewConsumerMiddleware(getConsumer()))
	kp.AddMiddleware(logmw.LogMiddleware())
	err := kp.Process(processUserLoggedInEvent)
	if err != nil {
		panic(err) // do better error handling
	}
}

func processUserLoggedInEvent(ctx context.Context, message *kafka.Message) error {
	// here, you can focus on your business logic.
	fmt.Printf("processing %v\n", message)
	time.Sleep(time.Millisecond * 200) // simulate long-running process
	return nil                         // or error
}

func getConfig() any {
	return nil // return your config
}
```
