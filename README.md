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

	processor := kp.NewKafkaProcessor("test", "retry-test", "dead-test", 10, "simple-service", kp.KafkaConfig{KafkaBootstrapServers: strings.Split(cfg.KafkaConfig.KafkaBootstrapServers, ",")})
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
