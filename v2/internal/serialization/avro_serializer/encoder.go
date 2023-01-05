package avro_serializer

import (
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/heetch/avro"
)

func (t serializer[MessageType]) Encode(message MessageType) ([]byte, error) {
	s := serde.BaseSerializer{}
	msgBytes, _, err := avro.Marshal(message)
	if err != nil {
		return nil, err
	}
	payload, err := s.WriteBytes(t.schemaID, msgBytes)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
