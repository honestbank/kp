---
sidebar_position: 5
---

# Logger Middleware
The zap logger middleware is a middleware for the KP framework that provides logging functionality using the Zap logger library. This middleware adds logging capabilities to your Kafka message processor, allowing you to log important information and errors during message processing.

### Usage
To use the zap logger middleware, you will need to import the package and create an instance of the Zap logger. Here's an example of how to set up the middleware:

```go
package main

import (
	"github.com/honestbank/kp/v2/middlewares/logger/zap"
	"go.uber.org/zap"
)

func main() {
	// Create a new Zap logger instance
	zapLogger, _ := zap.NewProduction()

	// Create the logger middleware
	loggerMiddleware := zap.NewLoggerMiddleware(zapLogger)

	// Add the logger middleware to your processor
	processor := kp.NewProcessor(
		kp.WithMiddleware(loggerMiddleware),
		// Add other options and middlewares
	)

	// Start the processor
	processor.Start()
}
```

In the above example, we create a new Zap logger instance using zap.NewProduction(). You can customize the logger configuration as per your requirements. Then, we create an instance of the zap logger middleware using zap.NewLoggerMiddleware(zapLogger). Finally, we add the logger middleware to our Kafka message processor using the kp.WithMiddleware option.

### Middleware Behavior
The zap logger middleware intercepts each message being processed by your processor and logs relevant information about the message, such as the key, offset, partition, and topic. It also logs the start and end of message processing, along with any errors that occur during processing.

The middleware uses the Zap logger's Info and Error levels to log the messages. The logging output can be customized by configuring the Zap logger instance.

### Example
```go
package main

func ExampleNewLoggerMiddleware() {
    // Create an initial logger for testing
    initialLogger := getInitialLogger()
    
    // Create the logger middleware
    middleware := zap.NewLoggerMiddleware(initialLogger)
    
    // Create a mock Kafka message
    topicName := "my-topic"
    mockKafkaMessage := &kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topicName, Offset: kafka.Offset(20), Partition: 2},
        Key:            []byte("my-key"),
    }
    
    // Process the mock Kafka message with the middleware
    _ = middleware.Process(context.Background(), mockKafkaMessage, func(ctx context.Context, item *kafka.Message) error {
        logger := zap.LoggerFromContext(ctx)
        logger.Info("processing message")
        logger.Error("returning an error")
        
        return errors.New("custom_error")
    })
// Output:
// {"level":"INFO","ts":"<timestamp>","msg":"Starting to process the message","key":"my-key","offset":"20","partition":2,"topic":"my-topic"}
// {"level":"INFO","ts":"<timestamp>","msg":"processing message","key":"my-key","offset":"20","partition":2,"topic":"my-topic"}
// {"level":"ERROR","ts":"<timestamp>","msg":"returning an error","key":"my-key","offset":"20","partition":2,"topic":"my-topic"}
// {"level":"ERROR","ts":"<timestamp>","msg":"Failed while processing the message","key":"my-key","offset":"20","partition":2,"topic":"my-topic","error":"custom_error"}
}
```

### FromContext
The zap logger middleware provides the LoggerFromContext function, which allows you to retrieve the configured logger from the context. This function can be used to obtain the logger instance with all the fields already set, enabling you to start logging without explicitly passing the logger around.

Here's an example of how you can use the LoggerFromContext function to retrieve the logger from the context and start logging:

```go
package main

func processMessage(ctx context.Context, item *kafka.Message) error {
	logger := zap.LoggerFromContext(ctx)
	logger.Info("Processing message")

	// Perform message processing logic

	return nil
}

```

## Conclusion

The zap logger middleware provides a seamless integration of the Zap logger library into your KP message processor. With this middleware, you can easily add logging capabilities to your processor, allowing you to monitor and debug message processing effectively.

It is recommended to customize the Zap logger configuration according to your application's logging requirements to get the desired logging output.
