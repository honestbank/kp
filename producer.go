package kp

import (
	"github.com/Shopify/sarama"

	"github.com/honestbank/kp/examples/simple/config"
)

var kpprodcer *Producer

type KPProducer interface {
	GetProducer() *Producer
	ProduceMessage(topic string, key string, message string) error
}

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer() KPProducer {
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	if kpprodcer == nil {
		saramaConfig := sarama.NewConfig()
		saramaConfig.Producer.Partitioner = sarama.NewRandomPartitioner
		saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
		saramaConfig.Producer.Return.Successes = true
		producer, err := sarama.NewSyncProducer([]string{config.KafkaConfig.KafkaBootstrap}, saramaConfig)
		if err != nil {
			panic(err)
		}
		return &Producer{
			producer: producer,
		}
	}
	return kpprodcer
}

func (p *Producer) ProduceMessage(topic string, key string, message string) error {
	messageToProduce := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
	}
	_, _, err := p.producer.SendMessage(messageToProduce)
	if err != nil {
		return err
	}
	return nil
}

func (p *Producer) GetProducer() *Producer {
	return p
}

func GetProducer() *Producer {
	if kpprodcer == nil {
		producer := NewProducer()
		kpprodcer = producer.GetProducer()
	}
	return kpprodcer
}
