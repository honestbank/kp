---
sidebar_position: 1
---

# Producer
KP exposes a producer that can produce avro formatted messages, this ensures backwards compatibility through schema registry.

:::tip
Initializing a producer also publishes the schema to schema registry, this allows you to fail the startup of your application if the schema has a breaking change.
:::

### Example {#example}
The following example sends a message to a Kafka topic if the type `UserLoggedIn` is backwards compatiable.

:::tip
Please check [this page](../introduction/configuration.md) for detailed configuration option
:::

```go
package main

import (
	"context"
	"github.com/honestbank/kp/v2/producer"
)

type UserLoggedIn struct {
	UserID string `avro:"user_id"`
}

func main() {
	p, err := producer.New[UserLoggedIn]("topic-name", getConfig())
	defer p.Flush()
	if err != nil {
		panic(err)
	}
	err = p.Produce(context.Background(), UserLoggedIn{})
	if err != nil {
		panic(err)
    }
}

func getConfig() any {
    return nil // return your config
}
```
