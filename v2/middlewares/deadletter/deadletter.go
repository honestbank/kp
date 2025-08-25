package deadletter

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/honestbank/kp/v2/internal/retrycounter"
	"github.com/honestbank/kp/v2/middlewares"
)

type Builder struct {
	onSuccess    func()
	onRetry      func(err error)
	onDeadLetter func(err error)
}

type deadletter struct {
	// required
	producer        Producer
	onProduceErrors func(err error)
	threshold       int
	// optional
	onSuccess    func()
	onRetry      func(err error)
	onDeadLetter func(err error)
}

type Producer interface {
	ProduceRaw(message *kafka.Message) error
}

func (r deadletter) Process(ctx context.Context, item *kafka.Message, next func(ctx context.Context, item *kafka.Message) error) error {
	err := next(ctx, item)
	if err == nil {
		r.onSuccess()
		return nil
	}
	if retrycounter.GetCount(item) < r.threshold {
		r.onRetry(err)
		return err
	}
	r.onDeadLetter(err)
	err = r.producer.ProduceRaw(&kafka.Message{Value: item.Value, Key: item.Key, Headers: item.Headers, Timestamp: item.Timestamp, TimestampType: item.TimestampType, Opaque: item.Opaque})
	if err != nil {
		r.onProduceErrors(err)
	}

	return nil
}

func NewDeadletterMiddleware(producer Producer, threshold int, onProduceErrors func(error)) middlewares.KPMiddleware[*kafka.Message] {
	return NewBuilder().Build(producer, threshold, onProduceErrors)
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

func (b *Builder) Build(producer Producer, threshold int, onProduceErrors func(error)) middlewares.KPMiddleware[*kafka.Message] {
	return deadletter{
		onProduceErrors: onProduceErrors,
		producer:        producer,
		threshold:       threshold,
		onSuccess:       b.onSuccess,
		onRetry:         b.onRetry,
		onDeadLetter:    b.onDeadLetter,
	}
}
