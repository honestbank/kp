---
sidebar_position: 3
---
# Configuration Options

KP doesn't expose a configuration option currently. This was done to minimize the API footprint.

But it does require a few things like kafka configuration and schema registry.
KP uses environment variable to configure connections.

## Kafka Configuration {#kafka-configuration}
The following environment variables are used for connecting to kafka:

- `KP_KAFKA_BOOTSTRAP_SERVERS` (required): URL to bootstrap servers eg: "localhost"
- `KP_SASL_MECHANISM` (optional): Learn about SASL authentication [here](https://docs.confluent.io/platform/current/kafka/authentication_sasl/index.html) (default: `""`)
- `KP_SECURITY_PROTOCOL` (optional): Learn more about auth [here](https://kafka.apache.org/25/javadoc/org/apache/kafka/common/security/auth/SecurityProtocol.html) (default: `""`)
- `KP_USERNAME` (optional): Value used for `sasl.username` (default: `""`).
- `KP_PASSWORD` (optional): Value used for `sasl.password` (default: `""`).
- `KP_SESSION_TIMEOUT_MS` (optional): Value used for `session.timeout.ms` and used by consumers (default: `6000`). 
- `KP_CONSUMER_AUTO_OFFSET_RESET` (optional): Value used for `auto.offset.reset` and used by consumers default: `"earliest"`

## Schema Registry Configuration {#schema-registry-configuration}
The following environment variable are used for schema registry.
As of now KP doesn't support running without a schema registry and there are no immediate active plans to change that.

- `KP_SCHEMA_REGISTRY_ENDPOINT` (required): URL to schema registry, eg: "http://localhost:8081"
- `KP_SCHEMA_REGISTRY_USERNAME` (optional): Username for basic authentication with schema registry (default: `""`)
- `KP_SCHEMA_REGISTRY_PASSWORD` (optional): Password for basic authentication with schema registry (default: `""`)
