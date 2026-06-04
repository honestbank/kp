package kp

import (
	"context"

	"github.com/IBM/sarama"
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

// producerMessageCarrier implements propagation.TextMapCarrier for sarama.ProducerMessage.
type producerMessageCarrier struct {
	msg *sarama.ProducerMessage
}

func (c producerMessageCarrier) Get(key string) string {
	for _, h := range c.msg.Headers {
		if string(h.Key) == key {
			return string(h.Value)
		}
	}

	return ""
}

func (c producerMessageCarrier) Set(key string, value string) {
	for i, h := range c.msg.Headers {
		if string(h.Key) == key {
			c.msg.Headers[i].Value = []byte(value)

			return
		}
	}
	c.msg.Headers = append(c.msg.Headers, sarama.RecordHeader{
		Key:   []byte(key),
		Value: []byte(value),
	})
}

func (c producerMessageCarrier) Keys() []string {
	keys := make([]string, len(c.msg.Headers))
	for i, h := range c.msg.Headers {
		keys[i] = string(h.Key)
	}

	return keys
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
	otel.GetTextMapPropagator().Inject(ctx, producerMessageCarrier{msg: messageToProduce})
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
