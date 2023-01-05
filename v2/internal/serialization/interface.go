package serialization

type Serializer[MessageType any] interface {
	Decode(payload []byte) (*MessageType, error)
	Encode(message MessageType) ([]byte, error)
}
