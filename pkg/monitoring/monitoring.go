package monitoring

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type OpenTelemetry interface {
	Start(ctx context.Context) (err error)
	Stop(ctx context.Context) (err error)
}

type openTelemetryImpl struct {
	serviceName   string
	environment   string
	endpoint      string
	resource      *resource.Resource
	traceExporter *otlptrace.Exporter
}

func NewOpenTelemetry(serviceName, environment, endpoint string) OpenTelemetry {
	ctx := context.Background()
	res, _ := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", environment),
		),
	)

	return &openTelemetryImpl{
		serviceName: serviceName,
		environment: environment,
		endpoint:    endpoint,
		resource:    res,
	}
}

// Start implements OpenTelemetry.
func (ot *openTelemetryImpl) Start(ctx context.Context) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(ot.endpoint),
		otlptracegrpc.WithDialOption(),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		otel.Handle(err)
		return
	}

	bsp := trace.NewBatchSpanProcessor(exporter)
	provider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(ot.resource),
		trace.WithSpanProcessor(bsp),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(provider)

	ot.traceExporter = exporter

	return
}

// Stop implements OpenTelemetry.
func (ot *openTelemetryImpl) Stop(ctx context.Context) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	err = ot.traceExporter.Shutdown(ctx)
	if err != nil {
		otel.Handle(err)
		return
	}

	return
}
