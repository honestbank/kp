package avro_serializer

import "github.com/honestbank/kp/v2/internal/serialization"

type serializer[MessageType any] struct {
	schemaID int
}

func New[MessageType any](schemaID int) serialization.Serializer[MessageType] {
	return serializer[MessageType]{
		schemaID: schemaID,
	}
}
