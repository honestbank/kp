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
		consumerGroup:   consumerGroup + "-" + "kp",
		backoffDuration: backoffDuration,
	}
}

func (k *KP) OnFailure(failure func(ctx context.Context, key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error) {
	k.onFailure = &failure
}

func (k *KP) Process(processor func(ctx context.Context, key string, message string, retries int, rawMessage *sarama.ConsumerMessage) error) {
	k.processor = processor
}

func (k *KP) Start(ctx context.Context) error {
	k.consumer = NewConsumer(k.topic, k.retryTopic, k.deadLetterTopic, k.retries, k.processor, k.onFailure, k.producer, k.backoffDuration)
	keepRunning := true
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
		if err = CreateConsumerSession(k, wg, ctx); err != nil {
			log.Println(err)
		}
	}() // need to handle the error in V2

	<-k.consumer.GetReady() // Await till the consumer has been set up

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			keepRunning = false
			k.Stop()
		case <-sigterm:
			keepRunning = false
			k.Stop()
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		return err
	}

	return nil
}

func (k *KP) Stop() {
	k.client.Close()
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
	} else {
		client.PauseAll()
	}

	*isPaused = !*isPaused
}

func CreateConsumerSession(k *KP, wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := k.client.Consume(ctx, []string{k.topic, k.retryTopic}, k.consumer); err != nil {
			return err
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return ctx.Err()
		}
		k.consumer.SetReady(make(chan bool))
	}
}
