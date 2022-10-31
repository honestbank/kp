package v2

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	backoff_policy "github.com/honestbank/backoff-policy"
	"github.com/honestbank/kp/v2/internal/consumer"
	"github.com/honestbank/kp/v2/internal/middleware"
	"github.com/honestbank/kp/v2/internal/retrycounter"
	"github.com/honestbank/kp/v2/internal/serialization"
	"github.com/honestbank/kp/v2/producer"
)

type kp[MessageType any] struct {
	topics           []string
	applicationName  string
	chain            middleware.Processor[*kafka.Message, error]
	retry            func(message *kafka.Message)
	sendToDeadLetter func(message *kafka.Message)
	cleanupCallbacks []func()
	backoff          backoff_policy.BackoffPolicy
	shouldContinue   bool
}

func (t *kp[MessageType]) WithRetryOrPanic(retryTopic string, retryCount int) KafkaProcessor[MessageType] {
	processor, err := t.WithRetry(retryTopic, retryCount)
	if err != nil {
		panic(err)
	}
	return processor
}

func (t *kp[MessageType]) WithDeadletterOrPanic(deadletterTopic string) KafkaProcessor[MessageType] {
	processor, err := t.WithDeadletter(deadletterTopic)
	if err != nil {
		panic(err)
	}
	return processor
}

func (t *kp[MessageType]) WithRetry(retryTopic string, retryCount int) (KafkaProcessor[MessageType], error) {
	t.topics = append(t.topics, retryTopic)
	p, err := producer.New[MessageType, int](retryTopic)
	if err != nil {
		return t, err
	}
	t.cleanupCallbacks = append(t.cleanupCallbacks, func() {
		p.Flush()
	})
	t.retry = func(message *kafka.Message) {
		c := retrycounter.GetCount(message)
		if c >= retryCount {
			t.sendToDeadLetter(message)

			return
		}
		retrycounter.SetCount(message, c+1)
		// increment message retry count
		// check if it should go to deadletter and send to either deadletter or retry topic
		_ = p.ProduceRaw(message)
	}

	return t, nil
}

func (t *kp[MessageType]) WithDeadletter(deadLetterTopic string) (KafkaProcessor[MessageType], error) {
	p, err := producer.New[MessageType, int](deadLetterTopic)
	if err != nil {
		return nil, err
	}
	t.cleanupCallbacks = append(t.cleanupCallbacks, func() {
		p.Flush()
	})
	t.sendToDeadLetter = func(message *kafka.Message) {
		_ = p.ProduceRaw(message)
	}

	return t, nil
}

func (t *kp[MessageType]) AddMiddleware(mw middleware.Middleware[*kafka.Message, error]) KafkaProcessor[MessageType] {
	t.chain.AddMw(mw)

	return t
}

func (t *kp[MessageType]) Stop() {
	t.shouldContinue = false
}

func (t *kp[MessageType]) Run(processor func(message MessageType) error) error {
	c, err := consumer.New(t.topics, t.applicationName)
	if err != nil {
		return err
	}
	t.chain.AddMw(middleware.FinalMw[*kafka.Message, error](func(msg *kafka.Message) error {
		message, err := serialization.Decode[MessageType](msg.Value)
		if err != nil {
			// do something with the err
		}
		err = processor(*message)
		if err != nil {
			t.retry(msg)
		}

		return err
	}))
	for t.shouldContinue {
		msg := c.GetMessage()
		if msg == nil {
			continue
		}
		t.chain.Process(msg)
	}
	for _, callback := range t.cleanupCallbacks {
		callback()
	}
	// need to clean up
	return nil
}

func New[MessageType any](topicName string, applicationName string) KafkaProcessor[MessageType] {
	return &kp[MessageType]{
		applicationName:  applicationName,
		chain:            middleware.New[*kafka.Message, error](),
		retry:            func(message *kafka.Message) {},
		sendToDeadLetter: func(message *kafka.Message) {},
		topics:           []string{topicName},
		backoff:          backoff_policy.NewExponentialBackoffPolicy(time.Millisecond*0, 0), // no backoff by default
		shouldContinue:   true,
	}
}
