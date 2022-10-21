package serialization

import (
	"github.com/heetch/avro"
)

func Decode[MessageType any](payload []byte) (*MessageType, error) {
	var emptyMessage MessageType
	typ, err := avro.TypeOf(emptyMessage)
	if err != nil {
		return nil, err
	}
	_, err = avro.Unmarshal(payload[5:], &emptyMessage, typ)
	if err != nil {
		return nil, err
	}

	return &emptyMessage, nil
}
