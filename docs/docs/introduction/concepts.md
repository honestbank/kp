---
sidebar_position: 4
---

# Kafka Concepts
A few of the Kafka concepts are abstracted away in this framework and this page attempts to clarify those concepts.

## Partitions {#partitions}
In Kafka, a topic is divided into one or more partitions, which allow the topic to scale horizontally. Each partition is an ordered, immutable sequence of messages that is continually appended to.

Producers can specify which partition to write to, allowing them to control the distribution of messages within a topic. If no partition is specified, the producer will use a default partitioning strategy, such as a round-robin approach.

Consumers can also specify which partitions to read from, allowing them to distribute the workload of reading and processing messages across multiple consumer instances.

## Offsets {#offsets}
Offsets are used to track the position of a consumer within a partition. Each time a consumer reads a message from a partition, it advances its offset by one. This allows the consumer to know which message to read next and enables Kafka to track the progress of the consumer.

Offsets can be committed by the consumer, allowing it to resume reading from the last committed offset if it is restarted or if there is a failure.

## Consumer Groups {#consumer-groups}
In Kafka, a consumer group is a set of consumer instances that work together to read and process messages from a topic. Each consumer in the group is assigned a set of partitions to read from, and the group works together to load balance the work of reading and processing the messages across all of its members.

If a consumer fails or is restarted, the consumer group will automatically reassign the partitions to another member to ensure that the messages are still being processed. This enables consumer groups to provide fault tolerance and allows for easy scaling of message processing.

## Replication {#replication}
Kafka allows for the replication of messages within a topic across multiple brokers. This ensures that the messages are highly available and can be recovered in the event of a broker failure.

Replication also allows for the use of multiple consumer groups to read from the same topic, allowing for parallel processing of the messages.

## Zookeeper {#zookeeper}
Zookeeper is a distributed coordination service that is used by Kafka to store metadata about the Kafka cluster and manage the distributed nature of the Kafka brokers.

Zookeeper is responsible for maintaining the state of the Kafka cluster, including the list of brokers, the topics that exist, the consumer groups that are active, and the consumer group assignments for each partition.

Kafka relies on Zookeeper to manage the distributed nature of the cluster and enable consumers and producers to discover and communicate with the brokers.

:::tip
KP disables autocommit and commits each message once they're finished processing.
This way we'll always wait for a message to complete processing (failing or succeeding) before committing.
:::
