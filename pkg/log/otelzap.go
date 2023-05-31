package log

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"strings"
)

type OtelZap struct {
	*zap.Logger
	ctx context.Context
}

func (l *OtelZap) Error(args ...interface{}) {
	l.Sugar().Error(l.otel(), args)
}
func (l *OtelZap) Errorf(template string, args ...interface{}) {
	l.Sugar().Error(l.otel(), fmt.Sprintf(template, args))
}
func (l *OtelZap) Info(args ...interface{}) {
	args = append(args, l.otel())
	l.Sugar().Info(args...)
}
func (l *OtelZap) Warn(args ...interface{}) {
	l.Sugar().Warn(l.otel(), args)
}
func (l *OtelZap) Panic(args ...interface{}) {
	l.Sugar().Panic(l.otel(), args)
}
func (l *OtelZap) Infof(template string, args ...interface{}) {
	args = append(args, l.otel())
	l.Sugar().Infof(template, args...)
}
func (l *OtelZap) Warnf(template string, args ...interface{}) {
	l.Warn(l.otel(), fmt.Sprintf(template, args))
}
func (l *OtelZap) Panicf(template string, args ...interface{}) {
	l.Sugar().Panic(l.otel(), fmt.Sprintf(template, args))
}
func (l *OtelZap) Debug(args ...interface{}) {
	l.Sugar().Debug(l.otel(), args)
}
func (l *OtelZap) Debugf(template string, args ...interface{}) {
	l.Sugar().Debug(l.otel(), fmt.Sprintf(template, args))
}

func (l *OtelZap) otel() string {
	span := trace.SpanFromContext(l.ctx)
	if !span.IsRecording() {
		return ""
	}
	builder := strings.Builder{}
	builder.WriteString("[")
	builder.WriteString(span.SpanContext().TraceID().String())
	builder.WriteString("]")
	return builder.String()
}
