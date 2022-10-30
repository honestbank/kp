package v2

type KafkaProcessor[MessageType any] interface {
	WithRetry(retryTopic string, retryCount int) (KafkaProcessor[MessageType], error)
	WithRetryOrPanic(retryTopic string, retryCount int) KafkaProcessor[MessageType]
	WithDeadletter(deadLetterTopic string) (KafkaProcessor[MessageType], error)
	WithDeadletterOrPanic(deadletterTopic string) KafkaProcessor[MessageType]
	Stop()
	Run(processor func(message MessageType) error) error
}
