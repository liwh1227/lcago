package log

import (
	"context"

	"go.uber.org/zap"
)

// wrapper.go 基于zap基础日志方法的封装

// Error zap log error
func Error(args ...interface{}) {
	getLogInstance().Error(args...)
}

// Errorf zap log error
func Errorf(template string, args ...interface{}) {
	getLogInstance().Errorf(template, args...)
}

// Info zap info
func Info(args ...interface{}) {
	getLogInstance().Info(args...)
}

// Infof zap infof
func Infof(template string, args ...interface{}) {
	//fmt.Sprintf(template, args)
	getLogInstance().Infof(template, args...)
}

// Warn zap warn
func Warn(args ...interface{}) {
	getLogInstance().Warn(args...)
}

// Warnf zap warn
func Warnf(template string, args ...interface{}) {
	getLogInstance().Warnf(template, args...)
}

// Debug zap debug
func Debug(args ...interface{}) {
	getLogInstance().Debug(args...)
}

// Debugf zap debug
func Debugf(template string, args ...interface{}) {
	getLogInstance().Debugf(template, args...)
}

// WithContext 带上下文的日志
func WithContext(ctx context.Context) *zap.SugaredLogger {
	log, ok := ctx.Value(getLogInstance().opts.CtxKey).(*zap.SugaredLogger)
	if ok {
		return log
	}
	return getLogInstance().SugaredLogger
}

// AddCtx 上下文中增加key-val
func AddCtx(ctx context.Context, field ...interface{}) (context.Context, *zap.SugaredLogger) {
	log := getLogInstance().SugaredLogger.With(field...)
	ctx = context.WithValue(ctx, l.opts.CtxKey, log)
	return ctx, log
}

// GetCtx 获取日志上下文
func GetCtx(ctx context.Context) *zap.SugaredLogger {
	log, ok := ctx.Value(getLogInstance().opts.CtxKey).(*zap.SugaredLogger)
	if ok {
		return log
	}
	return getLogInstance().SugaredLogger
}
