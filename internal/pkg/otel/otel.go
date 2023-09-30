package otel

import (
	"context"
	"flag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
)

func SetupOTelTracer(ctx context.Context) func(context.Context) error {
	flag.Parse()

	serviceName := os.Getenv("SERVICE_NAME")
	collectorURL := os.Getenv("EXPORTER_ENDPOINT")

	var logger = log.New(os.Stderr, serviceName, log.Ldate|log.Ltime|log.Llongfile)

	exporter, err := zipkin.New(
		collectorURL,
		zipkin.WithLogger(logger),
	)
	if err != nil {
		panic(err)
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown
}

func TraceIt(ctx context.Context, tracerName string) trace.Tracer {
	return otel.GetTracerProvider().Tracer(tracerName)
}
