package kp

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
)

type KP struct {
	consumer        ConsumerStruct
	topic           string
	retries         int
	deadLetterTopic string
	kafkaConfig     KafkaConfig
	producer        KPProducer
}

func NewKafkaProcessor(topic string, deadLetterTopic string, retries int, kafkaConfig KafkaConfig) KafkaProcessor {
	return &KP{
		topic:           topic,
		deadLetterTopic: deadLetterTopic,
		retries:         retries,
		kafkaConfig:     kafkaConfig,
		producer:        GetProducer(kafkaConfig),
	}
}

func (k *KP) Process(processor func(key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error) {
	k.consumer = NewConsumer(k.topic, k.deadLetterTopic, k.retries, processor, k.producer)
}

func (k *KP) Start() {
	group := "test"
	keepRunning := true
	log.Println("Starting a new Sarama consumer")
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true

	/**
	 * Setup a new Sarama consumer group
	 */

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(strings.Split(k.kafkaConfig.KafkaBootstrap, ","), group, saramaConfig)
	if err != nil {
		panic(err)
	}

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, []string{k.topic, k.deadLetterTopic}, &k.consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			k.consumer.SetReady(make(chan bool))
		}
	}()

	<-k.consumer.GetReady() // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}
