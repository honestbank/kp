package producer

import (
	"github.com/honestbank/kp/v2/config"
	"github.com/honestbank/kp/v2/internal/schemaregistry"
	"github.com/honestbank/kp/v2/internal/schemaregistry/avro_registry"
	"github.com/honestbank/kp/v2/internal/schemaregistry/noop_registry"
	"github.com/honestbank/kp/v2/internal/serialization"
	"github.com/honestbank/kp/v2/internal/serialization/avro_serializer"
	"github.com/honestbank/kp/v2/internal/serialization/json_serializer"
)

type Options[BodyType any] interface {
	PublishSchema() error
	Serialize(message BodyType) ([]byte, error)
}

type options[BodyType any] struct {
	registry   schemaregistry.SchemaRegistry[BodyType]
	serializer serialization.Serializer[BodyType]
}

func (o options[BodyType]) PublishSchema() error {
	_, err := o.registry.Publish()

	return err
}

func (o options[BodyType]) Serialize(message BodyType) ([]byte, error) {
	return o.serializer.Encode(message)
}

func WithJSONSerializer[BodyType any]() Options[BodyType] {
	return options[BodyType]{
		registry:   noop_registry.New[BodyType](),
		serializer: json_serializer.New[BodyType](),
	}
}

func defaultOptions[BodyType any](topicName string, cfg config.KPConfig) (Options[BodyType], error) {
	registry, err := avro_registry.New[BodyType](topicName, cfg.SchemaRegistryConfig)
	if err != nil {
		return nil, err
	}

	id, err := registry.Publish()
	if err != nil {
		return nil, err
	}

	serializer := avro_serializer.New[BodyType](*id)

	return options[BodyType]{
		registry:   registry,
		serializer: serializer,
	}, nil
}
