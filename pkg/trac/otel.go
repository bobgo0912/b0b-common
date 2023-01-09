package trac

import (
	"context"
	"fmt"
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/bobgo0912/b0b-common/pkg/constant"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

type OtelClient struct {
	Tp       *tracesdk.TracerProvider
	Tr       trace.Tracer
	ShutDown func(ctx context.Context) error
}

func NewOtelGrpc(ctx context.Context, options ...otlptracegrpc.Option) (*OtelClient, error) {
	client := otlptracegrpc.NewClient(options...)
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
	}
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(config.Cfg.ServiceName),
		semconv.ServiceVersionKey.String(config.Cfg.Version),
		attribute.String("env", string(config.Cfg.ENV)),
		attribute.String("nodeID", config.Cfg.NodeId),
	)
	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithSampler(GetSampler()),
		tracesdk.WithResource(resource),
	)
	otel.SetTracerProvider(tracerProvider)
	tracer := otel.GetTracerProvider().Tracer(
		config.Cfg.ServiceName,
		trace.WithInstrumentationVersion(config.Cfg.Version),
		trace.WithSchemaURL(semconv.SchemaURL),
	)
	return &OtelClient{Tp: tracerProvider, Tr: tracer, ShutDown: tracerProvider.Shutdown}, nil
}

func NewOtelHttp(ctx context.Context, options ...otlptracehttp.Option) (*OtelClient, error) {
	client := otlptracehttp.NewClient(options...)
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
	}
	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithSampler(GetSampler()),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.Cfg.ServiceName),
			semconv.ServiceVersionKey.String(config.Cfg.Version),
			attribute.String("env", string(config.Cfg.ENV)),
			attribute.String("nodeID", config.Cfg.NodeId),
		)),
	)
	otel.SetTracerProvider(tracerProvider)
	tracer := otel.GetTracerProvider().Tracer(
		config.Cfg.ServiceName,
		trace.WithInstrumentationVersion(config.Cfg.Version),
		trace.WithSchemaURL(semconv.SchemaURL),
	)

	return &OtelClient{Tp: tracerProvider, Tr: tracer, ShutDown: tracerProvider.Shutdown}, nil
}

func (o *OtelClient) StartTracer(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return o.Tr.Start(ctx, name, opts...)
}

func GetSampler() tracesdk.Sampler {
	switch config.Cfg.ENV {
	case constant.Dev:
		return tracesdk.AlwaysSample()
	case constant.Prod:
		return tracesdk.ParentBased(tracesdk.TraceIDRatioBased(0.5))
	default:
		return tracesdk.AlwaysSample()
	}
}
