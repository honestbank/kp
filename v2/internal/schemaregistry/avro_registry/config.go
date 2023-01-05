package avro_registry

type Config struct {
	Endpoint string `env:"KP_SCHEMA_REGISTRY_ENDPOINT" required:"true"`
	Username string `env:"KP_SCHEMA_REGISTRY_USERNAME" default:""`
	Password string `env:"KP_SCHEMA_REGISTRY_PASSWORD" default:""`
}
