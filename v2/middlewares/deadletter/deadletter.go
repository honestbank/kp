package deadletter

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/honestbank/kp/v2/internal/retrycounter"
	"github.com/honestbank/kp/v2/middlewares"
)

type Builder struct {
	onSuccess       func()          // called when the next middleware is executed successfully
	onRetry         func(err error) // called when the next middleware failed, but it didn't reach the maximum retry count
	onDeadLetter    func(err error) // called when the next middleware failed, and it reached the maximum retry count
	onProduceErrors func(err error) // called when it failed to produce the dead letter message
}

type deadletter struct {
	optionals Builder
	producer  Producer
	threshold int
}

type Producer interface {
	ProduceRaw(message *kafka.Message) error
}

func (r deadletter) Process(ctx context.Context, item *kafka.Message, next func(ctx context.Context, item *kafka.Message) error) error {
	err := next(ctx, item)
	if err == nil {
		r.optionals.onSuccess()
		return nil
	}
	if retrycounter.GetCount(item) < r.threshold {
		r.optionals.onRetry(err)
		return err
	}
	r.optionals.onDeadLetter(err)
	err = r.producer.ProduceRaw(&kafka.Message{Value: item.Value, Key: item.Key, Headers: item.Headers, Timestamp: item.Timestamp, TimestampType: item.TimestampType, Opaque: item.Opaque})
	if err != nil {
		r.optionals.onProduceErrors(err)
	}

	return nil
}

// Deprecated: use NewBuilder instead
func NewDeadletterMiddleware(producer Producer, threshold int, onProduceErrors func(error)) middlewares.KPMiddleware[*kafka.Message] {
	return NewBuilder().OnProduceErrors(onProduceErrors).Build(producer, threshold)
}

func NewBuilder() *Builder {
	return &Builder{
		onSuccess:    func() {},
		onRetry:      func(err error) {},
		onDeadLetter: func(err error) {},
	}
}

func (b *Builder) OnSuccess(f func()) *Builder {
	b.onSuccess = f

	return b
}

func (b *Builder) OnRetry(f func(err error)) *Builder {
	b.onRetry = f

	return b
}

func (b *Builder) OnDeadLetter(f func(err error)) *Builder {
	b.onDeadLetter = f

	return b
}

func (b *Builder) OnProduceErrors(f func(err error)) *Builder {
	b.onProduceErrors = f

	return b
}

func (b *Builder) Build(producer Producer, threshold int) middlewares.KPMiddleware[*kafka.Message] {
	return deadletter{
		producer:  producer,
		threshold: threshold,
		optionals: *b,
	}
}
