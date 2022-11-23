---
sidebar_position: 1
---

# Produce Raw Message
Producing a custom message without serialization

### Configuration {#configuration}

:::tip
Please check [this page](../introduction/configuration.md) for detailed configuration option
:::

The following example sends a message to a Kafka topic. This gives you full ability to fully customize the way you want to produce a message.

```go
package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/honestbank/kp/v2/producer"
)

type UserLoggedIn struct {
	UserID string `avro:"user_id"`
}

func main() {
	p, err := producer.New[UserLoggedIn]("topic-name", getConfig()) // this will still publish and check against schema registry.
	defer p.Flush()
	if err != nil {
		panic(err)
	}
	err = p.ProduceRaw(*kafka.Message{}) // producing raw will ignore the schema and won't be validated.
	if err != nil {
		panic(err)
	}
}

func getConfig() any {
	return nil // return your config
}
```