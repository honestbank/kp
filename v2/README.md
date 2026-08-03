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
| `PartitionAssignmentStrategy` | `partition.assignment.strategy` | unset → librdkafka default (`"range,roundrobin"`). Ignored when `GroupProtocol` is `"consumer"`. |
| `GroupProtocol` | `group.protocol` | unset → librdkafka default (`"classic"`). Set `"consumer"` for the KIP-848 protocol. |
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

## Rebalancing (classic vs. KIP-848 consumer protocol)

By default kp consumers use the **classic** rebalance protocol with librdkafka's `range,roundrobin` assignor — the **eager** protocol: on any rebalance (scale in/out, deploy), every consumer revokes all partitions and the whole group pauses until reassignment finishes.

For workloads that scale often (e.g. HPA), use the **KIP-848 consumer protocol** (`group.protocol=consumer`). Assignment is server-driven and incremental — only the partitions that actually move are paused — and it migrates online with a single rolling deploy.

```go
protocol := "consumer"
cfg := config.Kafka{
	BootstrapServers:  "...",
	ConsumerGroupName: "my-service",
	GroupProtocol:     &protocol,
}
```

Requirements and notes:

- **Brokers must support KIP-848** — Apache Kafka 4.0+ or Confluent Cloud.
- Under `group.protocol=consumer`, `session.timeout.ms`, `heartbeat.interval.ms` and `partition.assignment.strategy` are **server-managed**. kp automatically stops sending them (setting them client-side is a fatal librdkafka error), so `ConsumerSessionTimeoutMs` and `PartitionAssignmentStrategy` have no effect in this mode.
- **Migration is online.** Roll out `group.protocol=consumer` in a single deploy; when the first new pod joins, the group coordinator converts the group automatically and interoperates with any classic members still draining. No two-step dance, no downtime. Rolling back (removing the field) reverts the group once the last consumer-protocol member leaves.

> **Do not use `partition.assignment.strategy` for cooperative rebalancing on this stack.** librdkafka rejects a mixed eager+cooperative assignor list (`cooperative-sticky,range`) at startup, and a classic group has no zero-downtime path from eager to cooperative-sticky. Use the KIP-848 consumer protocol instead. `PartitionAssignmentStrategy` remains only for selecting a single classic assignor (e.g. `roundrobin`).
