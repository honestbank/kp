package kafkaheaders

import "github.com/confluentinc/confluent-kafka-go/kafka"

func Get(header string, message *kafka.Message) *string {
	for _, h := range message.Headers {
		if h.Key != header {
			continue
		}
		if h.Value == nil {
			return nil
		}
		s := string(h.Value)
		return &s
	}
	return nil
}
