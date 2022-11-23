---
sidebar_position: 1
---

# Getting Started with KP
KP is a Kafka message processing framework which makes retries, deadletters, tracing, measurements, etc easy.

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
A very basic example of processing messages would be something like the following

```go
package main

import (
	"context"
	"fmt"
	"time"

	v2 "github.com/honestbank/kp/v2"
)

type UserLoggedInEvent struct {
	UserID string
}

func main() {
	applicationName := "send-login-notification-worker"
	kp := v2.New[UserLoggedInEvent]("user-logged-in", getConfig())
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
By default, the above doesn't have any additional capacity, it's as simple as using a simple consumer provided by confluent.

KP however provides additional API to enable retries, deadletter topics etc.
