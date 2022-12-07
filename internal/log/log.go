package log

import (
	"github.com/bobgo0912/b0b-common/internal/constant"
	"go.uber.org/zap"
	"os"
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
)

func InitLog() {
	env := os.Getenv(constant.EnvName)
	if env == "prod" {
		logger, _ := zap.NewDevelopment()
		zap.ReplaceGlobals(logger)
	} else {
		logger, _ := zap.NewDevelopment()
		zap.ReplaceGlobals(logger)
	}
	Error = zap.S().Error
	Errorf = zap.S().Errorf
	Info = zap.S().Info
	Infof = zap.S().Infof
	Panic = zap.S().Panic
	Panicf = zap.S().Panicf
	Warn = zap.S().Warn
	Warnf = zap.S().Warnf
}
