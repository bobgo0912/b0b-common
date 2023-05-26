package middleware

import (
	"context"
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/bobgo0912/b0b-common/pkg/meter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"net/http"
	"sync"
	"time"
)

func Meter(next http.Handler) http.Handler {
	once := sync.Once{}
	once.Do(func() {
		meter.InitRequestCount()
		meter.InitRequestDuration()
		meter.InitRequestConcurrent()
	})
	meterm := otel.Meter(config.Cfg.ServiceName)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestStartTime := time.Now()
		ctx := r.Context()
		elapsedTime := float64(time.Since(requestStartTime)) / float64(time.Millisecond)

		// Record measurements
		attrs := semconv.HTTPServerMetricAttributesFromHTTPRequest(config.Cfg.ServiceName, r)
		value := attribute.String("path", r.URL.Path)
		attrs = append(attrs, value)

		meter.RequestCount.Add(ctx, 1, metric.WithAttributes(attrs...))
		meter.RequestDuration.Record(ctx, elapsedTime, metric.WithAttributes(attrs...))
		callback, _ := meterm.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {
			observer.ObserveInt64(meter.RequestConcurrent, -1)
			return nil
		}, meter.RequestConcurrent)
		next.ServeHTTP(w, r)
		callback.Unregister()
	})
}
