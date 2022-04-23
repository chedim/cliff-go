package core

import "go.uber.org/zap"

func LogDebug(args ...interface{}) {
	logger, _ := zap.NewDevelopment()
	l := logger.Sugar()
	l.Debug(args...)
}
