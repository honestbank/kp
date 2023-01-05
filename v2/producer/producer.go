package producer

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/honestbank/kp/v2/config"
	"github.com/honestbank/kp/v2/internal/tracing"
)

type producer[BodyType any] struct {
	schemaID int
	k        *kafka.Producer
	topic    string
	opts     Options[BodyType]
}

func (p producer[BodyType]) Produce(ctx context.Context, message BodyType) error {
	value, err := p.opts.Serialize(message)
	if err != nil {
		return err
	}

	partition := kafka.TopicPartition{
		Topic:     &p.topic,
		Partition: kafka.PartitionAny,
	}
	msg := &kafka.Message{
		TopicPartition: partition,
		Value:          value,
	}
	tracing.InjectTraceHeaders(ctx, msg)

	return p.ProduceRaw(msg)
}

func (p producer[BodyType]) Flush() error {
	// I just felt like 3000 because 3 second should be enough for messages to be flushed
	// we'll need to optimize this as we go
	p.k.Flush(3000)

	return nil
}

func (p producer[BodyType]) ProduceRaw(message *kafka.Message) error {
	// todo: maybe rename this method so that it tells people to not use it unless they know what they're doing.
	message.TopicPartition.Topic = &p.topic
	message.TopicPartition.Partition = kafka.PartitionAny
	return p.k.Produce(message, nil)
}

func New[MessageType any](topic string, cfg config.KPConfig) (Producer[MessageType], error) {
	opts, err := defaultOptions[MessageType](topic, cfg)
	if err != nil {
		return nil, err
	}

	return doNew[MessageType](topic, cfg, opts)
}

func NewWithOptions[MessageType any](topic string, cfg config.KPConfig, opts Options[MessageType]) (Producer[MessageType], error) {
	err := opts.PublishSchema()
	if err != nil {
		return nil, err
	}

	return doNew[MessageType](topic, cfg, opts)
}

func doNew[MessageType any](topic string, cfg config.KPConfig, opts Options[MessageType]) (Producer[MessageType], error) {
	k, err := kafka.NewProducer(config.GetKafkaConfig(cfg.KafkaConfig))
	if err != nil {
		return nil, err
	}

	return producer[MessageType]{
		k:     k,
		topic: topic,
		opts:  opts,
	}, nil
}
