package avro_registry

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/heetch/avro"

	"github.com/honestbank/kp/v2/config"
	internal "github.com/honestbank/kp/v2/internal/schemaregistry"
)

type client[MessageType any] struct {
	client    schemaregistry.Client
	topicName string
}

func getClient(cfg config.SchemaRegistry) (schemaregistry.Client, error) {
	client, err := schemaregistry.NewClient(schemaregistry.NewConfigWithAuthentication(
		cfg.Endpoint,
		cfg.Username,
		cfg.Password,
	))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c client[MessageType]) Publish() (*int, error) {
	var emptyMessage MessageType
	typ, err := avro.TypeOf(emptyMessage)
	if err != nil {
		return nil, err
	}
	schemaInfo := schemaregistry.SchemaInfo{
		Schema: typ.String(),
	}
	id, err := c.client.Register(fmt.Sprintf("%s-value", c.topicName), schemaInfo, true)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func New[MessageType any](topicName string, cfg config.SchemaRegistry) (internal.SchemaRegistry[MessageType], error) {
	c, err := getClient(cfg)

	if err != nil {
		return nil, err
	}

	return client[MessageType]{
		client:    c,
		topicName: topicName,
	}, nil
}
