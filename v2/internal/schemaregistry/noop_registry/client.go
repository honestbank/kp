package noop_registry

import "github.com/honestbank/kp/v2/internal/schemaregistry"

type client[BodyType any] struct{}

func (c client[BodyType]) Publish() (*int, error) {
	return nil, nil
}

func New[BodyType any]() schemaregistry.SchemaRegistry[BodyType] {
	return client[BodyType]{}
}
