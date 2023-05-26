package meter

import (
	"context"
	"github.com/bobgo0912/b0b-common/pkg/config"
	log "github.com/bobgo0912/b0b-common/pkg/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

type OtelMetricClient struct {
	Mp       *sdk.MeterProvider
	Mm       *metric.Meter
	ShutDown func(ctx context.Context) error
}

func NewOtelMetricHttp(ctx context.Context, opts ...otlpmetrichttp.Option) (*OtelMetricClient, error) {
	exporter, err := otlpmetrichttp.New(ctx, opts...)
	if err != nil {
		log.Error("NewOtelMeHttp fail err=", err)
		return nil, err
	}
	attributes := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(config.Cfg.ServiceName),
		semconv.ServiceVersionKey.String(config.Cfg.Version),
		attribute.String("env", string(config.Cfg.ENV)),
		attribute.String("nodeID", config.Cfg.NodeId),
	)

	// Pull-based Prometheus exporter
	prometheusExporter, err := prometheus.New()
	if err != nil {
		log.Error("NewOtelMeHttp fail err=", err)
		return nil, err
	}
	// Instantiate the OTLP HTTP exporter
	meterProvider := sdk.NewMeterProvider(
		sdk.WithResource(attributes),
		sdk.WithReader(sdk.NewPeriodicReader(exporter)),
		sdk.WithReader(prometheusExporter),
	)

	// Create an instance on a meter for the given instrumentation scope
	meter := meterProvider.Meter(
		"github.com/.../example/manual-instrumentation",
		metric.WithInstrumentationVersion("v0.0.0"),
	)
	otel.SetMeterProvider(meterProvider)

	return &OtelMetricClient{Mp: meterProvider, Mm: &meter, ShutDown: meterProvider.Shutdown}, nil
}
func NewOtelMetricGrpc(ctx context.Context, opts ...otlpmetricgrpc.Option) (*OtelMetricClient, error) {
	exporter, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		log.Error("NewOtelMetricGrpc fail err=", err)
		return nil, err
	}
	attributes := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(config.Cfg.ServiceName),
		semconv.ServiceVersionKey.String(config.Cfg.Version),
		attribute.String("env", string(config.Cfg.ENV)),
		attribute.String("nodeID", config.Cfg.NodeId),
	)
	// Pull-based Prometheus exporter
	prometheusExporter, err := prometheus.New()
	if err != nil {
		log.Error("NewOtelMeHttp fail err=", err)
		return nil, err
	}
	// Instantiate the OTLP HTTP exporter
	meterProvider := sdk.NewMeterProvider(
		sdk.WithResource(attributes),
		sdk.WithReader(sdk.NewPeriodicReader(exporter)),
		sdk.WithReader(prometheusExporter),
	)

	// Create an instance on a meter for the given instrumentation scope
	meter := meterProvider.Meter(
		config.Cfg.ServiceName,
		metric.WithInstrumentationVersion(config.Cfg.Version),
	)
	otel.SetMeterProvider(meterProvider)

	return &OtelMetricClient{Mp: meterProvider, Mm: &meter, ShutDown: meterProvider.Shutdown}, nil
}
