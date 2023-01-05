package json_serializer

import (
	"bytes"
	"encoding/json"

	"github.com/honestbank/kp/v2/internal/serialization"
)

type serializer[MessageType any] struct{}

func (s serializer[MessageType]) Decode(payload []byte) (*MessageType, error) {
	var m MessageType
	err := json.NewDecoder(bytes.NewReader(payload)).Decode(&m)
	if err != nil {
		return nil, err
	}

	return &m, err
}

func (s serializer[MessageType]) Encode(message MessageType) ([]byte, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(message)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), err
}

func New[MessageType any]() serialization.Serializer[MessageType] {
	return serializer[MessageType]{}
}
