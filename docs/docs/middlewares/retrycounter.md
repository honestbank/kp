---
sidebar_position: 7
---

# Retry Count
The retry count middleware is a middleware for the KP library that adds the current retry count to the context of the message processing. The retry count is the number of times the message has been retried by the retry middleware.

:::info
If you forget to add the middleware and tried to access the value, it's going to return `int(0)`
:::

### Example {#example}

To use the retry count middleware, you first need to add `RetryCountMiddleware`:

```go
package main

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/consumer"
	"github.com/honestbank/kp/v2/middlewares/retry_count"
)

func getConsumer() consumer.Consumer {
    return nil // omitted for brevity
}
func main() {
	kp := v2.New[kafka.Message]()
	kp.AddConsumer(getConsumer())
	kp.AddMiddleware(retry_count.NewRetryCountMiddleware()) // This adds retry count middleware
	err := kp.Process(processUserLoggedInEvent)
	if err != nil {
		panic(err) // do better error handling
	}
}

func processUserLoggedInEvent(ctx context.Context, message *kafka.Message) error {
	count := retry_count.FromContext(ctx) // this is how to get the count
	fmt.Printf("this message was processed %d times", count)
	return nil // or error
}
```
