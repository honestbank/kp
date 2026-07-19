# kp v2

Kafka message processing library built on [confluent-kafka-go](https://github.com/confluentinc/confluent-kafka-go) (librdkafka), with a middleware chain for retries, dead-lettering, tracing, and metrics. Messages are Avro-encoded against a schema registry.

```
go get github.com/honestbank/kp/v2
```

## Quick start

```go
package main

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	v2 "github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/config"
	kpconsumer "github.com/honestbank/kp/v2/consumer"
	"github.com/honestbank/kp/v2/middlewares/consumer"
	"github.com/honestbank/kp/v2/middlewares/gracefulshutdown"
)

func main() {
	cfg := config.Kafka{
		BootstrapServers:  "localhost",
		ConsumerGroupName: "my-service",
	}

	kafkaConsumer, err := kpconsumer.New([]string{"my-topic"}, cfg.WithDefaults())
	if err != nil {
		panic(err)
	}
	// Close after Run returns: sends LeaveGroup so the broker rebalances
	// immediately instead of waiting for session.timeout.ms.
	defer kafkaConsumer.Close()

	processor := v2.New[kafka.Message]()
	err = processor.
		AddMiddleware(gracefulshutdown.NewSignalMiddleware[*kafka.Message](processor.Stop)).
		AddMiddleware(consumer.NewConsumerMiddleware(kafkaConsumer)).
		Run(func(ctx context.Context, msg *kafka.Message) error {
			// process msg
			return nil
		})
	if err != nil {
		panic(err)
	}
}
```

Retry and dead-letter middlewares are available under `middlewares/retry` and `middlewares/deadletter`; see `kp_example_test.go` for a full chain.

## Configuration (`config.Kafka`)

| Field | librdkafka key | kp default (`WithDefaults()`) |
|---|---|---|
| `BootstrapServers` | `bootstrap.servers` | — (required) |
| `ConsumerGroupName` | `group.id` | — (required for consumers) |
| `SaslMechanism` | `sasl.mechanisms` | unset |
| `SecurityProtocol` | `security.protocol` | unset |
| `Username` / `Password` | `sasl.username` / `sasl.password` | unset |
| `ConsumerSessionTimeoutMs` | `session.timeout.ms` | `30000` |
| `ConsumerAutoOffsetReset` | `auto.offset.reset` | `"earliest"` |
| `ClientID` | `client.id` | `"rdkafka"` |
| `MaxMessageBytes` | `message.max.bytes` | unset |
| `SocketKeepaliveEnabled` | `socket.keepalive.enable` | `false` |
| `ConnectionsMaxIdleTimeoutMs` | `connections.max.idle.ms` | `30000` |
| `MaxPollIntervalMs` | `max.poll.interval.ms` | `30000` |
| `PartitionAssignmentStrategy` | `partition.assignment.strategy` | unset → librdkafka default (`"range,roundrobin"`) |
| `Debug` | `debug` | unset |

Pointer fields left `nil` (and not covered by a `WithDefaults()` value above) are never set on the librdkafka config, so librdkafka's own defaults apply.

Consumers always run with `enable.auto.commit=false`: each message is committed after your handler returns, giving at-least-once delivery.

> **Note:** kp's `MaxPollIntervalMs` default of `30000` is much stricter than librdkafka's own default (`300000`). A handler that takes longer than 30 s on a single message gets the consumer evicted from the group, triggering a rebalance and redelivery. Raise it if your handlers can be slow.

## Graceful shutdown

`processor.Stop()` only stops the polling loop. To leave the consumer group cleanly you must also close the consumer — otherwise the broker waits `session.timeout.ms` (30 s by default) before declaring the pod dead, and the pod's partitions stall for that long on every scale-in or deploy.

The pattern from the quick start handles this:

```
SIGTERM → Stop() → in-flight message finishes and commits
        → Run() returns → deferred Close() → LeaveGroup → immediate rebalance
```

On Kubernetes, set `terminationGracePeriodSeconds` greater than your worst-case single-message processing time.

## Partition assignment / cooperative rebalancing

By default (field unset) kp consumers use librdkafka's `range,roundrobin` — the **eager** protocol: on any rebalance (scale in/out, deploy), every consumer revokes all partitions and the whole group pauses until reassignment finishes.

For workloads that scale often (e.g. HPA), opt into **cooperative** rebalancing, where only the partitions that actually move are paused:

```go
strategy := "cooperative-sticky"
cfg := config.Kafka{
	BootstrapServers:            "...",
	ConsumerGroupName:           "my-service",
	PartitionAssignmentStrategy: &strategy,
}
```

### Migrating a live consumer group

A group cannot mix eager and cooperative members — switching directly causes `Inconsistent group protocol` errors during a rolling deploy. Migrate with two deploys:

1. Deploy with `"cooperative-sticky,range"` — the group keeps running eager, but every pod now also supports cooperative.
2. Deploy with `"cooperative-sticky"` — the group flips to cooperative.

Rolling back from cooperative to eager requires the same two steps in reverse.
