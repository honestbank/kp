---
sidebar_position: 8
---
# Retry
The retry middleware is a utility middleware that automatically retries processing a message in case of failure. When a message processing fails, the middleware produces the message to a retry topic. The middleware uses the `Producer` interface to produce the message, which means that the user can choose any Kafka client that implements this interface to produce the message.

:::tip
We send the message to the worker's retry topic to avoid delays in retries causing the consumer to lag.
:::

:::warning
Before you execute the Run method of KP, the retry topic must exist in Kafka. Failing to do so will produce the message to the retry topic, but it won't be attempted till next start of KP.
:::


The middleware also allows the user to specify a callback function, onProduceErrors, which is called in case of any error while producing the message. This can be useful for logging or handling any errors while producing the message.

To use the retry middleware, simply add it to the middleware chain and pass in a Producer implementation and the `onProduceErrors` callback function.

```go
retryMw := retry.NewRetryMiddleware(myProducer, func(err error) {
    fmt.Println(err)
})
kp.AddMiddleware(retryMw)
```

Note that the retry middleware should be added towards the end of the middleware chain, after all the processing middlewares. This is because the retry middleware should only be triggered in case of a processing failure and not for any other failure.

The retry middleware also increments a retry count on the message headers for every retry. This can be useful for tracking the number of times a message has been retried. The retry count can be retrieved using the `retry_count` middleware.

When a message exceeds the configured retry limit, it ignores the message by default. See [next page](./deadletter.md) to configure deadletters.
