package log

import (
	"context"
	"time"

	"github.com/zhangel/go-frame/config"
	"github.com/zhangel/go-frame/log/fields"
	"github.com/zhangel/go-frame/log/level"
	"github.com/zhangel/go-frame/log/logger"
)

func WithContext(ctx context.Context) logger.Logger {
	return DefaultLogger().WithContext(ctx)
}

func WithField(key string, value interface{}) logger.Logger {
	return DefaultLogger().WithField(key, value)
}

func WithFields(fields fields.Fields) logger.Logger {
	return DefaultLogger().WithFields(fields)
}

func WithConfig(config config.Config) logger.Logger {
	return DefaultLogger().WithConfig(config)
}

func WithLevel(level level.Level) logger.Logger {
	return DefaultLogger().WithLevel(level)
}

func WithRateLimit(interval time.Duration) logger.Logger {
	return DefaultLogger().WithRateLimit(interval)
}

func Tracef(format string, args ...interface{}) {
	DefaultLogger().Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	DefaultLogger().Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	DefaultLogger().Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	DefaultLogger().Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	DefaultLogger().Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	DefaultLogger().Fatalf(format, args...)
}

func Trace(args ...interface{}) {
	DefaultLogger().Trace(args...)
}

func Debug(args ...interface{}) {
	DefaultLogger().Debug(args...)
}

func Info(args ...interface{}) {
	DefaultLogger().Info(args...)
}

func Warn(args ...interface{}) {
	DefaultLogger().Warn(args...)
}

func Error(args ...interface{}) {
	DefaultLogger().Error(args...)
}

func Fatal(args ...interface{}) {
	DefaultLogger().Fatal(args...)
}

func LogLevel() level.Level {
	return DefaultLogger().LogLevel()
}
