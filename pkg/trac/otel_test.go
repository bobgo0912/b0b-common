package trac

import (
	"context"
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/bobgo0912/b0b-common/pkg/log"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"testing"
)

func TestNewOtelHttp(t *testing.T) {
	ctx, ca := context.WithCancel(context.Background())
	defer ca()
	log.InitLog()
	newConfig := config.NewConfig(config.Json)
	newConfig.Category = "../config"
	newConfig.InitConfig()
	config.Cfg.Version = "v0.1.0"
	otelHttp, err := NewOtelHttp(ctx, otlptracehttp.WithEndpoint("localhost:4318"), otlptracehttp.WithInsecure())
	defer otelHttp.ShutDown(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, span := otelHttp.StartTracer(ctx, "dtt")
	span.End()

}
func TestNewOtelGrpc(t *testing.T) {
	ctx, ca := context.WithCancel(context.Background())
	defer ca()
	log.InitLog()
	newConfig := config.NewConfig(config.Json)
	newConfig.Category = "../config"
	newConfig.InitConfig()
	config.Cfg.Version = "v0.1.0"
	otelGrpc, err := NewOtelGrpc(ctx, otlptracegrpc.WithEndpoint("localhost:4317"), otlptracegrpc.WithInsecure())
	defer otelGrpc.ShutDown(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, span := otelGrpc.StartTracer(ctx, "ggg")
	span.End()
}

func TestNewOtelGrpcJaeger(t *testing.T) {
	ctx, ca := context.WithCancel(context.Background())
	defer ca()
	log.InitLog()
	newConfig := config.NewConfig(config.Json)
	newConfig.Category = "../config"
	newConfig.InitConfig()
	config.Cfg.Version = "v0.1.0"
	otelGrpc, err := NewOtelGrpcJaeger(jaeger.WithEndpoint("http://localhost:14268/api/traces"))
	defer otelGrpc.ShutDown(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, span := otelGrpc.StartTracer(ctx, "ggg")
	span.End()

	c := make(chan struct{})
	<-c
}
