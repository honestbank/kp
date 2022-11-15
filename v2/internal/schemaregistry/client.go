package schemaregistry

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/heetch/avro"

	"github.com/honestbank/kp/v2/config"
)

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

func Publish[MessageType any](topicName string, cfg config.SchemaRegistry) (*int, error) {
	client, err := getClient(cfg)
	if err != nil {
		return nil, err
	}
	var emptyMessage MessageType
	typ, err := avro.TypeOf(emptyMessage)
	if err != nil {
		return nil, err
	}
	schemaInfo := schemaregistry.SchemaInfo{
		Schema: typ.String(),
	}
	id, err := client.Register(fmt.Sprintf("%s-value", topicName), schemaInfo, true)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
