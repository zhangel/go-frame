package log

import (
	"context"
	"errors"
	"time"

	"github.com/zhangel/go-frame.git/config"
	"github.com/zhangel/go-frame.git/log/fields"
	"github.com/zhangel/go-frame.git/log/level"
	"github.com/zhangel/go-frame.git/log/logger"
)

type _Logger struct {
	opts Options
}

func NewLogger(opts ...Option) (logger.Logger, error) {
	log := &_Logger{}
	for _, o := range opts {
		if err := o(&log.opts); err != nil {
			return nil, err
		}
	}

	if log.opts.encoder == nil {
		return nil, errors.New("no log encoder available")
	}

	if log.opts.writer == nil {
		return nil, errors.New("no log writer available")
	}

	logger := log.WithConfig(config.GlobalConfig())
	if log.opts.minLevelOp != nil {
		logger = log.opts.minLevelOp(logger)
	}

	return logger, nil
}

func (s *_Logger) WithContext(ctx context.Context) logger.Logger {
	return newEntry(s).WithContext(ctx)
}

func (s *_Logger) WithField(key string, value interface{}) logger.Logger {
	return newEntry(s).WithField(key, value)
}

func (s *_Logger) WithFields(fields fields.Fields) logger.Logger {
	return newEntry(s).WithFields(fields)
}

func (s *_Logger) WithConfig(config config.Config) logger.Logger {
	return newEntry(s).WithConfig(config)
}

func (s *_Logger) WithLevel(level level.Level) logger.Logger {
	return newEntry(s).WithLevel(level)
}

func (s *_Logger) WithRateLimit(interval time.Duration) logger.Logger {
	return newEntry(s).WithRateLimit(interval)
}

func (s *_Logger) WithFuncName(enable bool) logger.Logger {
	return newEntry(s).WithFuncName(enable)
}

func (s *_Logger) Tracef(format string, args ...interface{}) {
	newEntry(s).Tracef(format, args...)
}

func (s *_Logger) Debugf(format string, args ...interface{}) {
	newEntry(s).Debugf(format, args...)
}

func (s *_Logger) Infof(format string, args ...interface{}) {
	newEntry(s).Infof(format, args...)
}

func (s *_Logger) Warnf(format string, args ...interface{}) {
	newEntry(s).Warnf(format, args...)
}

func (s *_Logger) Errorf(format string, args ...interface{}) {
	newEntry(s).Errorf(format, args...)
}

func (s *_Logger) Fatalf(format string, args ...interface{}) {
	newEntry(s).Fatalf(format, args...)
}

func (s *_Logger) Trace(args ...interface{}) {
	newEntry(s).Trace(args...)
}

func (s *_Logger) Debug(args ...interface{}) {
	newEntry(s).Debug(args...)
}

func (s *_Logger) Info(args ...interface{}) {
	newEntry(s).Info(args...)
}

func (s *_Logger) Warn(args ...interface{}) {
	newEntry(s).Warn(args...)
}

func (s *_Logger) Error(args ...interface{}) {
	newEntry(s).Error(args...)
}

func (s *_Logger) Fatal(args ...interface{}) {
	newEntry(s).Fatal(args...)
}

func (s *_Logger) Close() error {
	if s.opts.writer != nil {
		_ = s.opts.writer.Close()
	}
	return nil
}
