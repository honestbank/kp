---
sidebar_position: 4
---

# Kafka Concepts
A few of the Kafka concepts are abstracted away in this framework and this page attempts to clarify those concepts.

### Message Bus {#message-bus}
While a message queue allows possibility to delete a message or re-queue a message, message bus is a append only messaging system.
Kafka uses offsets to track pointers of the consumers in a partition.
Because of committing nature of Kafka, there's no way to mark a message as failed so that we can retry later.

### Producers {#producers}
Producers are systems that can write messages into a Kafka topic.
When a message is produced to a topic, Kafka appends this to a partition at the end.

### Consumers {#consumers}
Consumers are systems that read messages from a topic.
They read messages of the partition they're assigned to, and they read them in order.
Each time a consumer commits a message, Kafka keeps track of that so that consumer can resume where they left off.

### Commits {#commits}
When a consumer reads a chunk of message (or 1 message at a time) it can choose to commit. Think of this like a checkpoint.
Once it's committed, this consumer if restarted will resume processing from that offset.

Kafka has an autocommit functionality where a processor can continue to read the message.
Most client can automatically commit offsets in a certain interval, but that can be disabled.
If processing messages only once is a goal, we can disable autocommit and manually commit instead.

:::tip
KP disables autocommit and commits each message once they're finished processing.
This way we'll always wait for a message to complete processing (failing or succeeding) before committing.
:::

