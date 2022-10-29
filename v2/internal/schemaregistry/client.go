package schemaregistry

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/heetch/avro"

	"github.com/honestbank/kp/v2/internal/config"
)

func getClient() (schemaregistry.Client, error) {
	cfg, err := config.LoadConfig[Config]()
	if err != nil {
		return nil, err
	}
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

func Publish[MessageType any](topicName string) (*int, error) {
	client, err := getClient()
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
