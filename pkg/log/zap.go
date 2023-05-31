package log

import (
	"go.uber.org/zap"
)

type ZapLog struct {
	Zap *zap.Logger
}

func (l *ZapLog) Error(args ...interface{}) {
	l.Zap.Sugar().Error(args)
}
func (l *ZapLog) Errorf(template string, args ...interface{}) {
	l.Zap.Sugar().Errorf(template, args)
}
func (l *ZapLog) Info(args ...interface{}) {
	l.Zap.Sugar().Info(args)
}
func (l *ZapLog) Warn(args ...interface{}) {
	l.Zap.Sugar().Warn(args)
}
func (l *ZapLog) Panic(args ...interface{}) {
	l.Zap.Sugar().Panic(args)
}
func (l *ZapLog) Infof(template string, args ...interface{}) {
	l.Zap.Sugar().Infof(template, args)
}
func (l *ZapLog) Warnf(template string, args ...interface{}) {
	l.Zap.Sugar().Warnf(template, args)
}
func (l *ZapLog) Panicf(template string, args ...interface{}) {
	l.Zap.Sugar().Panicf(template, args)
}
func (l *ZapLog) Debug(args ...interface{}) {
	l.Zap.Sugar().Debug(args)
}
func (l *ZapLog) Debugf(template string, args ...interface{}) {
	l.Zap.Sugar().Debugf(template, args)
}
