---
sidebar_position: 10
---
# graceful shutdown
The `gracefulshutdown` middleware is designed to handle the `syscall.SIGINT`, `syscall.SIGTERM`, `syscall.SIGHUP`, and `syscall.SIGQUIT` signals gracefully, by allowing you to stop your processor in a controlled manner.

The `gracefulshutdown` middleware simply calls the next function passed to it, and does not modify or affect the item in any way. However, it does register a callback function that will be called when a SIGINT signal is received by the process.

This allows you to cleanly shut down your processor, for example, by closing any open resources or waiting for in-flight messages to be processed before exiting.

## Usage {#usage}
To use the sigint middleware, you first need to create a new instance of it by calling `gracefulshutdown.NewSigIntMiddleware[T](stopFunction func())`. The stopFunction parameter is a callback function that will be called when a one of the `syscall.SIGINT`, `syscall.SIGTERM`, `syscall.SIGHUP`, or `syscall.SIGQUIT` signal is received.

Then, you can use the middleware in your pipeline by passing it to the Process method of KP.

Here's an example of how you might use this middleware in a simple pipeline:

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	v2 "github.com/honestbank/kp/v2"
	consumer2 "github.com/honestbank/kp/v2/consumer"
	"github.com/honestbank/kp/v2/middlewares/consumer"
	"github.com/honestbank/kp/v2/middlewares/gracefulshutdown"
)

func main() {
	c, err := consumer.New([]string{"topic-1", "topic-2"}, getConfig())
	if err != nil {
		// handle err
	}
	consumerMiddleware := consumer.NewConsumerMiddleware(c)
	processor := v2.New[kafka.Message]()
	processor.AddMiddleware(gracefulshutdown.NewSignalMiddleware(func() {
		processor.Stop()
	}))
	processor.AddMiddleware(consumerMiddleware)
	processor.Process(processUserLoggedInEvent)
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
This will stop the processor when the process receives one of `syscall.SIGINT`, `syscall.SIGTERM`, `syscall.SIGHUP`, or `syscall.SIGQUIT` signal.
