package schemaregistry

type SchemaRegistry[BodyType any] interface {
	Publish() (*int, error)
}
