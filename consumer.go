package kp

import (
	"log"

	"github.com/Shopify/sarama"
)

type Consumer interface {
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	ConsumeClaim(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error
}

// ConsumerStruct represents a Sarama consumer group consumer
type ConsumerStruct struct {
	topic           string
	deadLetterTopic string
	ready           chan bool
	processor       func(message string) error
	producer        KPProducer
	retries         int
}

func NewConsumer(topic string, deadLetterTopic string, retries int, processor func(message string) error, producer KPProducer) ConsumerStruct {
	return ConsumerStruct{
		ready:           make(chan bool),
		processor:       processor,
		producer:        producer,
		topic:           topic,
		deadLetterTopic: deadLetterTopic,
		retries:         retries,
	}
}

func (Consumer *ConsumerStruct) GetReady() chan bool {
	return Consumer.ready
}

func (consumer *ConsumerStruct) SetReady(ready chan bool) {
	consumer.ready = ready
}

func (consumer *ConsumerStruct) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *ConsumerStruct) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *ConsumerStruct) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message := <-claim.Messages():
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			unmarshaledMessage, retries, err := UnmarshalStringMessage(string(message.Value))
			if err != nil {
				log.Printf("Error unmarshaling message: %v", err)
			}
			if retries >= consumer.retries {
				log.Printf("Message has exceeded retries, sending to dead letter topic")
				session.MarkMessage(message, "")
				return nil
			}
			err = consumer.processor(unmarshaledMessage)
			if err != nil {

				marshaledMessage := MarshalStringMessage(unmarshaledMessage, retries+1)
				err = consumer.producer.ProduceMessage(consumer.deadLetterTopic, string(message.Key), marshaledMessage)
				if err != nil {
					log.Println("ERROR OCCURED")
				}
			}
			session.MarkMessage(message, "")

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
