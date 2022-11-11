---
sidebar_position: 3
---

# Producer
KP exposes a producer which can produce avro formatted messages ensuring backwards compatibility through schema registry.

:::tip
Initializing producer also publishes the schema to schema registry, this allows you to fail the startup of your application If schema has a breaking change.
:::

### Example {#example}
The following example sends a message to a kafka topic.

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
	p, err := producer.New[UserLoggedIn]("topic-name")
	defer p.Flush()
	if err != nil {
		panic(err)
	}
	err = p.Produce(context.Background(), UserLoggedIn{})
	if err != nil {
		panic(err)
    }
}
```
