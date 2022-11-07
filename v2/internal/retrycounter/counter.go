package retrycounter

import (
	"strconv"

	"github.com/honestbank/kp/v2/internal/kafkaheaders"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const retryCountHeader = "x-retry-count"

func GetCount(message *kafka.Message) int {
	val := kafkaheaders.Get(retryCountHeader, message)
	if val == nil {
		return 0
	}
	count, err := strconv.Atoi(*val)
	if err != nil {
		return 0
	}
	return count
}

func SetCount(message *kafka.Message, count int) {
	kafkaheaders.Set(message, retryCountHeader, strconv.Itoa(count))
}
