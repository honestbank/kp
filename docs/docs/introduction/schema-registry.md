---
sidebar_position: 5
---
# Schema Registry
Confluent kafka comes with a schema registry, and you can choose to deploy this schema registry by yourself.
We can use a docker image to deploy schema registry on premises as well.

Schema Registry enables kafka clients to publish the message schema for a topic.
It allows schema evolution with certain conditions. One of the example is backwards compatibility.

:::tip
Schema registry is made available once you sign up to confluent cloud, or you can choose to host the entire stack yourself.
See `v2/examples/full` for a minimal example.
:::

## KP and Schema Registry {#kp-schema-registry}
As of now, kp requires a schema registry to make sure we only deploy backwards compatible applications.
We also have no immediate plans to change this any time soon.
If you genuinely feel kp would help more without a schema registry, please reach out to us.
