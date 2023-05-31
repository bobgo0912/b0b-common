package log

import (
	"context"
	"github.com/bobgo0912/b0b-common/pkg/constant"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

type LogFather interface {
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Panic(args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
}

var (
	Log        LogFather
	once       sync.Once
	otelLogger *otelzap.Logger
	gzap       *zap.Logger
)

var (
	Error  func(args ...interface{})
	Errorf func(template string, args ...interface{})
	Info   func(args ...interface{})
	Warn   func(args ...interface{})
	Panic  func(args ...interface{})
	Infof  func(template string, args ...interface{})
	Warnf  func(template string, args ...interface{})
	Panicf func(template string, args ...interface{})
	Debug  func(args ...interface{})
	Debugf func(template string, args ...interface{})
)

func InitLog() error {
	env := os.Getenv(constant.EnvName)
	var globalLogger *zap.Logger
	if env == "prod" {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return err
		}
		globalLogger = logger
		zap.ReplaceGlobals(logger)
	} else {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return err
		}
		globalLogger = logger
		zap.ReplaceGlobals(logger)
	}
	gzap = globalLogger
	Error = zap.S().Error
	Errorf = zap.S().Errorf
	Info = zap.S().Info
	Infof = zap.S().Infof
	Panic = zap.S().Panic
	Panicf = zap.S().Panicf
	Warn = zap.S().Warn
	Warnf = zap.S().Warnf
	Debugf = zap.S().Debugf
	Debug = zap.S().Debug
	log := &ZapLog{Zap: globalLogger}
	Log = log
	return nil
}

//func Otel(ctx context.Context) *OtelZap {
//	once.Do(func() {
//		if gzap == nil {
//			InitLog()
//		}
//		otelLogger = otelzap.New(gzap)
//		otelzap.ReplaceGlobals(otelLogger)
//	})
//	return &OtelZap{
//		Logger: gzap,
//		ctx:    ctx,
//	}
//}

func Otel(ctx context.Context) otelzap.LoggerWithCtx {
	once.Do(func() {
		if gzap == nil {
			InitLog()
		}
		otelLogger = otelzap.New(gzap, otelzap.WithMinLevel(zapcore.InfoLevel), otelzap.WithTraceIDField(true))
		otelzap.ReplaceGlobals(otelLogger)
	})
	return otelLogger.Ctx(ctx)
}
