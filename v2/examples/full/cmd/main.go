package main

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	"github.com/honestbank/kp_demo/server"
	"github.com/honestbank/kp_demo/worker"
)

func main() {
	defer setupTracing()()
	if len(os.Args) <= 1 {
		panic("needs a second argument")
	}
	switch os.Args[1] {
	case "server":
		server.StartServer()
	case "worker":
		worker.Process()
	default:
		panic("must select server or worker")
	}
}

func setupTracing() func() {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(os.Getenv("COLLECTOR_URL"))))
	if err != nil {
		panic(err)
	}
	resc := resource.Default()
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resc),
		trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("form-service-web"))),
	)
	otel.SetTracerProvider(tp)
	return func() {
		exporter.Shutdown(context.Background())
		tp.Shutdown(context.Background())
	}
}
