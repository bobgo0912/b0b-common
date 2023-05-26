package meter

import (
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/bobgo0912/b0b-common/pkg/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var RequestCount metric.Int64Counter
var RequestDuration metric.Float64Histogram
var RequestConcurrent metric.Int64ObservableGauge

func InitRequestCount() (metric.Int64Counter, error) {
	meter := otel.Meter(config.Cfg.ServiceName)
	histogram, err := meter.Int64Counter("request_count",
		metric.WithDescription("Incoming request count"),
		metric.WithUnit("request"))
	if err != nil {
		log.Error("InitRequestCount fail err=", err.Error())
		return nil, err
	}
	RequestCount = histogram

	return histogram, err
}
func InitRequestDuration() (metric.Float64Histogram, error) {
	meter := otel.Meter(config.Cfg.ServiceName)
	histogram, err := meter.Float64Histogram("duration",
		metric.WithDescription("Incoming end to end duration"),
		metric.WithUnit("ms"))
	if err != nil {
		log.Error("InitRequestDuration fail err=", err.Error())
		return nil, err
	}
	RequestDuration = histogram
	return histogram, err
}
func InitRequestConcurrent() (metric.Int64ObservableGauge, error) {
	meter := otel.Meter(config.Cfg.ServiceName)
	histogram, err := meter.Int64ObservableGauge("concurrent",
		metric.WithDescription("concurrent req"))
	if err != nil {
		log.Error("InitRequestConcurrent fail err=", err.Error())
		return nil, err
	}
	RequestConcurrent = histogram
	return histogram, err
}
