package trac

import (
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"net/http"
)

const otelName = "b0b-common/http"

// TraceSpan is a middleware that initialize a tracing span and injects span
// context to r.Context(). In one word, this middleware kept an eye on the
// whole HTTP request that the server receives.
func TraceSpan(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer(otelName).Start(r.Context(), "TraceSpan")
		if span == nil {
			// Tracer not found, just skip.
			next.ServeHTTP(w, r)
		}
		defer span.End()
		span.SetAttributes(semconv.HTTPSchemeHTTP)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
