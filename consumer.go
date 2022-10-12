package kp

import (
	"context"
	"errors"
	"time"

	"github.com/Shopify/sarama"

	backoff_policy "github.com/honestbank/backoff-policy"
)

type KPConsumer interface {
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	ConsumeClaim(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error
	GetReady() chan bool
	SetReady(chan bool)
	Process(ctx context.Context, message *sarama.ConsumerMessage) error
}

// ConsumerStruct represents a Sarama consumer group consumer
type ConsumerStruct struct {
	topic           string
	deadLetterTopic string
	retryTopic      string
	ready           chan bool
	Processor       func(ctx context.Context, key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error
	producer        KPProducer
	retries         int
	backoffPolicy   backoff_policy.BackoffPolicy
	onFailure       *func(ctx context.Context, key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error
}

func NewConsumer(topic string, retryTopic string, deadLetterTopic string, retries int, processor func(ctx context.Context, key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error, onFailure *func(ctx context.Context, key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error, producer KPProducer, backoffPolicyTime time.Duration) KPConsumer {
	return &ConsumerStruct{
		ready:           make(chan bool),
		Processor:       processor,
		producer:        producer,
		topic:           topic,
		deadLetterTopic: deadLetterTopic,
		retries:         retries,
		retryTopic:      retryTopic,
		backoffPolicy:   backoff_policy.NewExponentialBackoffPolicy(backoffPolicyTime, retries),
		onFailure:       onFailure,
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

func (consumer *ConsumerStruct) Process(ctx context.Context, message *sarama.ConsumerMessage) error {
	if message == nil {
		return errors.New("error while trying to consume nil message")
	}
	unmarshalMessage, retries, err := UnmarshalStringMessage(string(message.Value))
	if err != nil {
		return err
	}
	if retries >= consumer.retries {
		err = consumer.producer.ProduceMessage(ctx, consumer.deadLetterTopic, string(message.Key), unmarshalMessage)
		if consumer.onFailure != nil {
			err = (*consumer.onFailure)(ctx, string(message.Key), unmarshalMessage, retries, message)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}

		return nil
	}
	consumer.backoffPolicy.Execute(func(marker backoff_policy.Marker) {
		err := consumer.Processor(ctx, string(message.Key), unmarshalMessage, retries, message)
		if err != nil {
			marker.MarkFailure()
			if err != nil {
				marshaledMessage := MarshalStringMessage(unmarshalMessage, retries+1)
				err = consumer.producer.ProduceMessage(ctx, consumer.retryTopic, string(message.Key), marshaledMessage)
				if err != nil {
					return //need to handle the error in V2
				}
			}

			return
		}
		marker.MarkSuccess()
	})

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
			err := consumer.Process(session.Context(), message)
			if err != nil {
				return err
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
