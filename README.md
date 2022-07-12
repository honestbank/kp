# kp
Library for handing Kafka messages with retries

## How it works
Each KP instance create a kafka sarama client which will listen to 2 topics: main topic and retry topic. The Retry
topic will be prepended with the consumer group name such that each instance could have its own retry and dead letter
topics. This will enable a case such that you have two services that listen to the same topic but one service might
fail the topic, with this it would only retry on that service and not the other one.

Code Example: (see [examples](https://github.com/honestbank/kp/tree/main/examples))

```golang
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	processor := kp.NewKafkaProcessor("test", "retry-test", "dead-test", 10, "simple-service", kp.KafkaConfig{KafkaBootstrapServers: strings.Split(cfg.KafkaConfig.KafkaBootstrapServers, ",")}, 0)
	// retry topic becomes "simple-service-retry-test"
	// dead letter topic becomes "simple-service-dead-test"
	processor.Process(func(key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error {
		if message == "fail" {
			return errors.New("failed")
		}
		log.Println("message content:" + message)
		return nil
	})

	processor.Start()

}
```

To disable retries, just set retries to 0

## Exponential backoff

This library supports exponential backoff. To use a backoff policy, set the backoffDuration to anything above 1. To
disable the backoff policy, set the backoffDuration to 0. Backoff happens on the whole client which will slow down all
mesaages. When a message is successful, the backoffDuration is reduced.

## Topics

KP Uses 3 topics: main topic, retry topic and dead letter topic.
You are able to modify the retry and dead letter topics like so:

```golang

processor := kp.NewKafkaProcessor(
	cfg.Kafka.CardReplacementTopic,
	"retry-"+cfg.Kafka.CardReplacementTopic,
	"deadletter-"+cfg.Kafka.CardReplacementTopic,
	10,
	"card-delivery-service",
	kp.KafkaConfig{KafkaBootstrapServers: strings.Split(cfg.Kafka.KafkaBootstrap, ",")},
	time.Second*5,
	)
```

This is useful when you have a worker that listens to multiple topics and each topic needs to have its own retry and
dead letter topic.

For retry and deadletter topics they are generated with the consumer group.
ex:
```
  topic: test
  retrytopic: retry
  deadlettertopic: deadletter
  consumergroup: group

resulting topics:
    retrytopic: group-retry
    deadlettertopic: group-deadletter


resulting topic format: <consumer group>-<retry|deadletter topic>
```
