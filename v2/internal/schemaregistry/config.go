package schemaregistry

type Config struct {
	Endpoint string `env:"SCHEMA_REGISTRY_ENDPOINT" required:"true"`
	Username string `env:"SCHEMA_REGISTRY_USERNAME" default:""`
	Password string `env:"SCHEMA_REGISTRY_PASSWORD" default:""`
}
