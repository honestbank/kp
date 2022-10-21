package serialization

import (
	"reflect"

	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/heetch/avro"
)

func Encode(message any, schemaID int) ([]byte, error) {
	if message == nil {
		return nil, nil
	}
	s := serde.BaseSerializer{}
	val := reflect.ValueOf(message)
	if val.Kind() == reflect.Ptr {
		// avro.TypeOf expects an interface containing a non-pointer
		message = val.Elem().Interface()
	}
	msgBytes, _, err := avro.Marshal(message)
	if err != nil {
		return nil, err
	}
	payload, err := s.WriteBytes(schemaID, msgBytes)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
