package logger

import (
	"context"
	"time"

	"github.com/zhangel/go-frame/config"
	"github.com/zhangel/go-frame/log/fields"
	"github.com/zhangel/go-frame/log/level"
)

type Logger interface {
	WithContext(ctx context.Context) Logger
	WithField(key string, value interface{}) Logger
	WithFields(fields fields.Fields) Logger
	WithConfig(config config.Config) Logger
	WithLevel(minLevel level.Level) Logger
	WithRateLimit(interval time.Duration) Logger
	WithFuncName(enable bool) Logger

	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	LogLevel() level.Level

	Close() error
}

type LoggerFactory func() (Logger, error)
