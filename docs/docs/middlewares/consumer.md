---
sidebar_position: 2
---
# Consumer
The consumer middleware is responsible for populating the chain with messages from the Kafka topic(s). Consumer middleware retrieves a message from the Kafka topic using the Consumer interface and passes it to the next middleware in the chain.

The consumer middleware should be added before any other middleware that needs a message to be present.

:::tip
By default, the middleware chain is empty. If the consumer of the library does not wish to read from Kafka, they can skip using the consumer middleware and even implement alternative message brokers, such as Amazon SQS.
:::

```go
package main

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	v2 "github.com/honestbank/kp/v2"
	"time"

	consumer2 "github.com/honestbank/kp/v2/consumer"
	"github.com/honestbank/kp/v2/middlewares/consumer"
)

func main() {
	c, err := consumer.New([]string{"topic-1", "topic-2"}, getConfig())
	if err != nil {
		// handle err
	}
	consumerMiddleware := consumer.NewConsumerMiddleware(c)
	processor := v2.New[kafka.Message]()
	processor.AddMiddleware(consumerMiddleware)
	processor.Process(processUserLoggedInEvent)
}

func processUserLoggedInEvent(ctx context.Context, message *kafka.Message) error {
	// here, you can focus on your business logic.
	fmt.Printf("processing %v\n", message)
	time.Sleep(time.Millisecond * 200) // simulate long running process
	return nil                         // or error
}

func getConfig() any {
	return nil // return your config
}
```