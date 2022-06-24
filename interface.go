package kp

type KafkaProcessor interface {
	Process(processor func(message string) error)
	Start()
}
