package kafkaheaders

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

func Set(message *kafka.Message, header string, value string) {
	for i, h := range message.Headers {
		if h.Key != header {
			continue
		}
		message.Headers[i] = kafka.Header{Key: header, Value: []byte(value)}
		return
	}
	message.Headers = append(message.Headers, kafka.Header{Key: header, Value: []byte(value)})
}
