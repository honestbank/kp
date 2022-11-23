---
sidebar_position: 3
---
# Configuration Options

KP can be configured using a `config.KPConfig` object.

Here's the structure of KPConfig struct
```go
package config

type KPConfig struct {
	KafkaConfig          Kafka
	SchemaRegistryConfig SchemaRegistry
}

type Kafka struct {
	ConsumerGroupName        string
	BootstrapServers         string
	SaslMechanism            *string
	SecurityProtocol         *string
	Username                 *string
	Password                 *string
	ConsumerSessionTimeoutMs *int
	ConsumerAutoOffsetReset  *string
}

type SchemaRegistry struct {
    Endpoint string
    Username string
    Password string
}
```

## Kafka Configuration {#kafka-configuration}
The following fields are used to configure KP

- `BootstrapServers` (required): URL to bootstrap servers eg: "localhost"
- `SaslMechanism` (optional): Learn about SASL authentication [here](https://docs.confluent.io/platform/current/kafka/authentication_sasl/index.html) (default: `""`)
- `SecurityProtocol` (optional): Learn more about auth [here](https://kafka.apache.org/25/javadoc/org/apache/kafka/common/security/auth/SecurityProtocol.html) (default: `""`)
- `Username` (optional): Value used for `sasl.username` (default: `""`).
- `Password` (optional): Value used for `sasl.password` (default: `""`).
- `ConsumerGroupName` (required, consumer only): To configure consumer group id
- `ConsumerSessionTimeoutMs` (optional, consumer only): Value used for `session.timeout.ms` (default: `6000`). 
- `ConsumerAutoOffsetReset` (optional, consumer only): Value used for `auto.offset.reset` (default: `"earliest"`).

## Schema Registry Configuration {#schema-registry-configuration}
The following fields are used for schema registry.
As of now KP doesn't support running without a schema registry and there are no immediate active plans to change that.

- `Endpoint` (required): URL to schema registry, eg: "http://localhost:8081"
- `Username` (optional): Username for basic authentication with schema registry (default: `""`)
- `Password` (optional): Password for basic authentication with schema registry (default: `""`)
