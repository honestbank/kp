package producer

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/honestbank/kp/v2/config"
	"github.com/honestbank/kp/v2/internal/middleware"
)

type untypedProducer struct {
	producer           *kafka.Producer
	topic              string
	middlewarePipeline middleware.Processor[*kafka.Message, error]
}

func (u untypedProducer) Produce(ctx context.Context, message *kafka.Message) error {
	return u.middlewarePipeline.Process(ctx, message)
}

func (u untypedProducer) SetMiddlewares(middlewares []middleware.Middleware[*kafka.Message, error]) {
	pipeline := middleware.New[*kafka.Message, error]()
	for _, m := range middlewares {
		pipeline.AddMiddleware(m)
	}
	pipeline.AddMiddleware(middleware.FinalMiddleware(func(ctx context.Context, message *kafka.Message) error {
		message.TopicPartition.Topic = &u.topic
		message.TopicPartition.Partition = kafka.PartitionAny

		return u.producer.Produce(message, nil)
	}))
	u.middlewarePipeline = pipeline
}

func (u untypedProducer) Flush() error {
	// I just felt like 3000 because 3 second should be enough for messages to be flushed
	// we'll need to optimize this as we go
	u.producer.Flush(3_000)

	return nil
}

func NewUntyped(topic string, cfg config.Kafka) (UntypedProducer, error) {
	p, err := kafka.NewProducer(config.GetKafkaConfig(cfg))
	if err != nil {
		return nil, err
	}
	pipeline := middleware.New[*kafka.Message, error]()

	pipeline.AddMiddleware(middleware.FinalMiddleware[*kafka.Message, error](func(ctx context.Context, item *kafka.Message) error {
		item.TopicPartition.Topic = &topic
		item.TopicPartition.Partition = kafka.PartitionAny

		return p.Produce(item, nil)
	}))

	return untypedProducer{
		producer:           p,
		topic:              topic,
		middlewarePipeline: pipeline,
	}, nil
}
