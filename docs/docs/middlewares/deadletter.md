---
sidebar_position: 9
---
# Deadletter
The deadletter middleware is used to send messages that have failed processing to a deadletter topic after a certain number of retries. This can be useful for identifying and debugging messages that are consistently causing processing errors.

To use the deadletter middleware, you will need to provide a producer interface that can be used to send messages to the deadletter topic. You will also need to specify a threshold for the number of retries before a message is sent to the deadletter topic. Optionally, you can also provide a function that will be called whenever there is an error while producing a message to the deadletter topic.

To use the deadletter middleware, you will need to provide a producer interface that can be used to send messages to the deadletter topic. You will also need to specify a threshold for the number of retries before a message is sent to the deadletter topic. Optionally, you can also provide a function that will be called whenever there is an error while producing a message to the deadletter topic.

:::info
While you absolutely can use retries without deadletters, it'll probably be hard to setup re-processing of failed items. Using deadletters is highly recommended.
:::

```go
deadletterMiddleware := middlewares.NewDeadletterMiddleware(myProducer, 5, myOnProduceErrorsFunc)
processor := kp.AddMiddleware(deadletterMiddleware)
```

The deadletter middleware should be placed after the retry middleware in the middleware chain to ensure that messages are properly retried before being sent to the deadletter topic.

:::warning
It is important to note that messages sent to the deadletter topic are not deleted from Kafka, so you will need to implement a separate process for cleaning up these messages to prevent them from accumulating over time.
:::
