---
sidebar_position: 4
---
# Metrics
The metrics middleware is a middleware for the KP library that tracks various metrics about the message processing, such as the duration of each operation and the number of successes and failures. These metrics are collected using the Prometheus library and can be pushed to a Prometheus gateway to be stored and visualized.

:::warning
The metrics middleware should be added after the backoff middleware to avoid measuring the wait times. However, you can choose to add it before the backoff middleware if you want to include the wait times in the metrics.
:::

The metrics middleware is a useful tool for tracking the performance and stability of the KP library and the message processing. By adding it to the middleware chain, you can collect and monitor various metrics about the message processing and use them to improve the system.

### Metrics Collected
The metrics middleware collects the following metrics:

`kp_operation_duration_milliseconds`: A histogram of the duration of each operation, in milliseconds, labeled by the result (success or failure) and the error (if any).

### Example {#example}

To use the metrics middleware, you first need to import the middlewares/measurements package and call the `NewMeasurementMiddleware` function:

The gatewayURL argument is the URL of the Prometheus gateway to which the metrics should be pushed, and the applicationName argument is the name of the application that will be used as the job name in the metrics.

Then, you can add the metrics middleware to the KP instance by calling the AddMiddleware method:

```go
package main

import (
	"context"
	"fmt"
	consumer2 "github.com/honestbank/kp/v2/consumer"
	"github.com/honestbank/kp/v2/middlewares/consumer"
	"time"

	backoff_policy "github.com/honestbank/backoff-policy"
	"github.com/honestbank/backoff-policy/policies"
	v2 "github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/middlewares/measurement"
)

type UserLoggedInEvent struct {
	UserID string
}

func getConsumer() consumer2.Consumer {
    // omitted for brevity
	return nil
}

func main() {
	applicationName := "send-login-notification-worker"
	kp := v2.New[kafka.Message]()
	kp.AddMiddleware(consumer.NewConsumerMiddleware(getConsumer()))
	kp.AddMiddleware(measurement.NewMeasurementMiddleware("http://path/to/push/gateway", applicationName)) // simply add a measurement middleware to get free metrics
	err := kp.Process(processUserLoggedInEvent)
	if err != nil {
		panic(err) // do better error handling
	}
}

func processUserLoggedInEvent(ctx context.Context, message *kafka.Message) error {
	// here, you can focus on your business logic.
	fmt.Printf("processing %v\n", message)
	time.Sleep(time.Millisecond * 200) // simulate long running process
	return nil                         // or error
}

func getConfig() any {
	return nil // return your config
}
```

:::warning
The metrics middleware starts a background job that pushes the metrics to the Prometheus gateway every 5 seconds. Make sure that the gateway is running and accessible from the application.
:::
