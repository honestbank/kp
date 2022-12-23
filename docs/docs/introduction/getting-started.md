---
sidebar_position: 1
---

# Getting Started with KP
In this guide, we will walk through the steps to set up and use the KP library to process messages from a Kafka topic. We will start by installing the library and its dependencies, then we will write a custom processor function and add middlewares to customize the behavior of the message processing.

It helps you write applications without worrying too much about the retries, deadlettering, backing off when things go wrong, etc. As a developer, simply focus on business logic and business logic only.
Your code will be free of retries, backoffs and even tracing, but you'll still able to get all those for free.

It also comes with a opinionated producer which can produce typed messages into a given topic.

:::tip
Please check [this page](../introduction/concepts.md) for very brief concepts of Kafka.
:::


### Installation {#installation}
Adding KP as dependency on your project is as simple as running

```bash
go get github.com/honestbank/kp/v2
```

### Requirements {#requirements}
To use KP, you'll need the following:
- go version 1.18 is needed because of the usage of generics
- Schema registry to register schema of your topics (can't be disabled currently, maybe in future)
- `cgo` needs to be enabled during build of your go code.

:::tip
Please check [this page](../introduction/configuration.md) for detailed configuration option
:::

## Basic Example (publishing messages) {#basic-example-producer}
A producer can be initialized with a type, and it'll automatically publish the type to schema registry.
It requires some environment variables to be set.

```go
package main

import (
	"context"

	"github.com/honestbank/kp/v2/producer"
)

type UserLoggedInEvent struct {
	UserID string
}

func main() {
	p, err := producer.New[UserLoggedInEvent]("user-logged-in", getConfig())
	if err != nil {
		panic(err) // do better error handling
	}
	defer p.Flush()
	err = p.Produce(context.Background(), UserLoggedInEvent{UserID: "dummy"})
	if err != nil {
		panic(err)
    }
}

func getConfig() any {
    return nil // return your config
}
```

## Basic Example (processing messages) {#basic-example-kp}
### Processor function
The processor function is the main function that is called for each message in the middleware chain. It should perform the main processing of the message and return an error if any.

### Adding Middlewares
To add middlewares to the KP instance, you can use the AddMiddleware method. This method takes a middleware as an argument and returns a new KP instance with the middleware added to the chain.

### Running the Message Processing
To start the message processing, you can call the `Run` method on the KP instance and pass the processor function as an argument.

### Stopping the Message Processing
To stop the message processing, you can call the Stop method on the KP instance:

This method stops the message processing and releases resources. It is important to call the Stop method when you are finished using the KP instance.

Here's what all the above might look like:

```go
package main

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"time"

	v2 "github.com/honestbank/kp/v2"
)

type UserLoggedInEvent struct {
	UserID string
}

func main() {
	kp := v2.New[kafka.Message]()
	kp.AddMiddleware(getConsumer())
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
