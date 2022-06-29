package kp

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
)

type KP struct {
	consumer        KPConsumer
	topic           string
	retryTopic      string
	consumerGroup   string
	retries         int
	deadLetterTopic string
	kafkaConfig     KafkaConfig
	producer        KPProducer
	client          sarama.ConsumerGroup
	backoffDuration time.Duration
	processor       func(ctx context.Context, key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error
	onFailure       *func(ctx context.Context, key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error
}

func NewKafkaProcessor(topic string, retryTopic string, deadLetterTopic string, retries int, consumerGroup string, kafkaConfig KafkaConfig, backoffDuration time.Duration) KafkaProcessor {
	return &KP{
		topic:           topic,
		deadLetterTopic: consumerGroup + "-" + deadLetterTopic,
		retryTopic:      consumerGroup + "-" + retryTopic,
		retries:         retries,
		kafkaConfig:     kafkaConfig,
		producer:        GetProducer(kafkaConfig),
		consumerGroup:   consumerGroup,
		backoffDuration: backoffDuration,
	}
}

func (k *KP) OnFailure(failure func(ctx context.Context, key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error) {
	k.onFailure = &failure
}

func (k *KP) Process(processor func(ctx context.Context, key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error) {
	k.processor = processor
}

func (k *KP) Start(ctx context.Context) {
	k.consumer = NewConsumer(k.topic, k.retryTopic, k.deadLetterTopic, k.retries, k.processor, k.onFailure, k.producer, k.backoffDuration)
	keepRunning := true
	log.Println("Starting a new Sarama consumer")
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true

	/**
	 * Setup a new Sarama consumer group
	 */

	ctx, cancel := context.WithCancel(ctx)
	client, err := sarama.NewConsumerGroup(k.kafkaConfig.KafkaBootstrapServers, k.consumerGroup, saramaConfig)
	k.client = client
	if err != nil {
		panic(err)
	}

	otelsarama.WrapConsumerGroupHandler(k.consumer)

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := k.client.Consume(ctx, []string{k.topic, k.retryTopic}, k.consumer); err != nil {
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

func (k *KP) Stop() {
	log.Println("Stopping Kafka consumer")
	k.client.Close()
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
