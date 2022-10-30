package retrycounter

import (
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const retryCountHeader = "x-retry-count"

func GetCount(message *kafka.Message) int {
	for _, header := range message.Headers {
		if header.Key == retryCountHeader {
			count, err := strconv.Atoi(string(header.Value))
			if err != nil {
				return 0
			}
			return count
		}
	}
	return 0
}

func SetCount(message *kafka.Message, count int) {
	for _, header := range message.Headers {
		if header.Key == retryCountHeader {
			header.Value = []byte(strconv.Itoa(count))
			return
		}
	}
	message.Headers = append(message.Headers, kafka.Header{Key: retryCountHeader, Value: []byte(strconv.Itoa(count))})
}
