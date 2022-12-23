---
sidebar_position: 2
---

# Consumer
The Consumer interface provides a way to retrieve and commit messages from a Kafka topic. It has two methods:

- `GetMessage() *kafka.Message`: This method retrieves a message from the Kafka topic. It returns a pointer to a kafka.Message struct, which contains the message value, key, and other metadata.
- `Commit(message *kafka.Message) error`: This method commits a message to the Kafka topic. It takes a pointer to a kafka.Message struct as an argument, and returns an error if the commit fails.



:::tip
While it is possible to use the Consumer interface directly, it is intended to be used through the consumer middleware
:::

### Example {#example}
Here is an example of how to use the Consumer interface to retrieve and commit messages from a Kafka topic:

:::tip
Please check [configuration page](../introduction/configuration.md) for detailed configuration option
:::

```go
package main

import (
	"github.com/honestbank/kp/v2/consumer"
)

func main() {
	c, err := consumer.New([]string{"topic-1", "topic-2"}, getConfig())
	if err != nil {
		// handle err
    }
	msg := c.GetMessage()
	if msg == nil {
		// handle error or no message available
	}

	// Process the message
	// ...

	// Commit the message to the Kafka topic
	if err := c.Commit(msg); err != nil {
		// handle error
	}
}

func getConfig() any {
	return nil // return your config
}
```
### Notes {#notes}
- If the `GetMessage()` method returns nil, it could mean that there are no messages available or that an error occurred, generally errors are recovered internally by confluent client.
- The `Commit()` method should only be called after you have successfully processed the message and are ready to commit it to the Kafka topic.
