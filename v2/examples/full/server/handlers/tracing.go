package handlers

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func Tracing(handler http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		traceCTX := propagation.TraceContext{}
		request = request.WithContext(traceCTX.Extract(request.Context(), propagation.HeaderCarrier(request.Header)))
		tracer := otel.GetTracerProvider().Tracer("http_handler")
		ctx, span := tracer.Start(request.Context(), "request_received")
		propagation.TraceContext{}.Inject(ctx, propagation.HeaderCarrier(writer.Header()))
		handler.ServeHTTP(writer, request.WithContext(ctx))
		span.End()
	}
}
