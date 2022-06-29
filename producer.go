package kp

import (
	"context"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
)

var kpprodcer *Producer

type KPProducer interface {
	GetProducer(kafkaConfig KafkaConfig) *Producer
	ProduceMessage(ctx context.Context, topic string, key string, message string) error
}

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(kafkaConfig KafkaConfig) KPProducer {
	if kpprodcer == nil {
		saramaConfig := sarama.NewConfig()
		saramaConfig.Producer.Partitioner = sarama.NewRandomPartitioner
		saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
		saramaConfig.Producer.Return.Successes = true
		producer, err := sarama.NewSyncProducer(kafkaConfig.KafkaBootstrapServers, saramaConfig)
		if err != nil {
			panic(err)
		}
		producer = otelsarama.WrapSyncProducer(saramaConfig, producer)

		return &Producer{
			producer: producer,
		}
	}

	return kpprodcer
}

func (p *Producer) ProduceMessage(ctx context.Context, topic string, key string, message string) error {
	messageToProduce := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
	}
	otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(messageToProduce))
	_, _, err := p.producer.SendMessage(messageToProduce)
	if err != nil {
		return err
	}

	return nil
}

func (p *Producer) GetProducer(kafkaConfig KafkaConfig) *Producer {
	return p
}

func GetProducer(kafkaConfig KafkaConfig) *Producer {
	if kpprodcer == nil {
		producer := NewProducer(kafkaConfig)
		kpprodcer = producer.GetProducer(kafkaConfig)
	}

	return kpprodcer
}
