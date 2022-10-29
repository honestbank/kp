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
	// first byte is the magic byte and the next 4 bytes are the schema ID (int) which is why we start reading from the 5th position
	_, err = avro.Unmarshal(payload[5:], &emptyMessage, typ)
	if err != nil {
		return nil, err
	}

	return &emptyMessage, nil
}
