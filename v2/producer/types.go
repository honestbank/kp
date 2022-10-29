package producer

type Producer[BodyType any, KeyType KeyTypes] interface {
	Flush() error
	Produce(message KafkaMessage[BodyType, KeyType]) error
}

type KafkaMessage[BodyType any, KeyType any] struct {
	Body BodyType
	Key  *KeyType
}

type KeyTypes interface {
	int | int64 | string
}
