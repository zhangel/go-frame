package log

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zhangel/go-frame.git/config"
	"github.com/zhangel/go-frame.git/config/watcher"
	"github.com/zhangel/go-frame.git/lifecycle"
	"github.com/zhangel/go-frame.git/log/entry"
	"github.com/zhangel/go-frame.git/log/fields"
	"github.com/zhangel/go-frame.git/log/level"
	"github.com/zhangel/go-frame.git/log/logger"
)

type rateLimitKey struct {
	file string
	line int
}

var (
	lastTimeMap = make(map[rateLimitKey]time.Time)
	lastTimeMu  sync.Mutex
)

type _Entry struct {
	logger          *_Logger
	ctx             context.Context
	fields          *fields.Fields
	config          *config.Config
	minLevel        *uint32
	interval        *time.Duration
	funcNameTracing *bool
	cancel          []func()
}

func newEntry(logger *_Logger) *_Entry {
	minLevel := uint32(level.TraceLevel)
	interval := time.Duration(0)
	fields := fields.Fields{}
	funcNameTracing := false
	return &_Entry{
		logger:          logger,
		fields:          &fields,
		minLevel:        &minLevel,
		interval:        &interval,
		funcNameTracing: &funcNameTracing,
	}
}

func (s *_Entry) WithContext(ctx context.Context) logger.Logger {
	e := newEntry(s.logger)
	e.ctx = ctx
	e.fields = s.fields
	e.config = s.config
	e.minLevel = s.minLevel
	e.interval = s.interval
	e.funcNameTracing = s.funcNameTracing

	return e
}

func (s *_Entry) WithField(key string, value interface{}) logger.Logger {
	return s.WithFields(fields.Fields{fields.Field{K: key, V: value}})
}

func (s *_Entry) WithFields(fs fields.Fields) logger.Logger {
	e := newEntry(s.logger)
	f := fields.Fields{}
	f = append(f, *s.fields...)
	f = append(f, fs...)
	e.fields = &f
	e.ctx = s.ctx
	e.config = s.config
	e.minLevel = s.minLevel
	e.interval = s.interval
	e.funcNameTracing = s.funcNameTracing

	return e
}

func (s *_Entry) WithConfig(config config.Config) logger.Logger {
	if config == nil {
		return s
	}

	e := newEntry(s.logger)
	e.fields = s.fields
	e.ctx = s.ctx
	e.config = &config
	e.interval = s.interval
	e.funcNameTracing = s.funcNameTracing

	if l, err := level.ParseLevel(config.String(flagMinLevel)); err != nil {
		minLevelU32 := uint32(level.TraceLevel)
		e.minLevel = &minLevelU32
	} else {
		minLevelU32 := uint32(l)
		e.minLevel = &minLevelU32
	}

	interval := config.Duration(flagRateLimit)
	if interval != 0 {
		e.interval = &interval
	}

	s.cancel = append(s.cancel, config.Watch(watcher.NewHelper(flagMinLevel, func(v string, deleted bool) {
		if deleted {
			atomic.StoreUint32(e.minLevel, uint32(level.TraceLevel))
		} else if l, err := level.ParseLevel(v); err == nil {
			atomic.StoreUint32(e.minLevel, uint32(l))
		} else {
			return
		}
	})))

	s.cancel = append(s.cancel, config.Watch(watcher.NewHelper(flagRateLimit, func(v string, deleted bool) {
		interval := config.Duration(flagRateLimit)
		if interval != 0 {
			e.interval = &interval
		}
	})))

	return e
}

func (s *_Entry) WithLevel(level level.Level) logger.Logger {
	e := newEntry(s.logger)
	e.fields = s.fields
	e.ctx = s.ctx
	e.config = s.config
	minLevelU32 := uint32(level)
	e.minLevel = &minLevelU32
	e.interval = s.interval
	e.funcNameTracing = s.funcNameTracing

	return e
}

func (s *_Entry) WithRateLimit(interval time.Duration) logger.Logger {
	e := newEntry(s.logger)
	e.fields = s.fields
	e.ctx = s.ctx
	e.config = s.config
	e.minLevel = s.minLevel
	e.interval = &interval
	e.funcNameTracing = s.funcNameTracing

	return e
}

func (s *_Entry) WithFuncName(enable bool) logger.Logger {
	e := newEntry(s.logger)
	e.fields = s.fields
	e.ctx = s.ctx
	e.config = s.config
	e.minLevel = s.minLevel
	e.interval = s.interval
	e.funcNameTracing = &enable

	return e
}

func (s *_Entry) log(level level.Level, a ...interface{}) {
	e := entry.Entry{}
	e.Time = time.Now()
	e.Caller = GetCaller(true)
	if s.isRateLimited(e.Caller) {
		return
	}
	e.Level = level
	e.Msg = fmt.Sprint(a...)
	e.Fields = *s.fields

	if s.ctx != nil {
		for i := range ctxFieldProviders {
			e.Fields = append(e.Fields, ctxFieldProviders[i](s.ctx)...)
		}
	}

	if s.logger.opts.encoderExt != nil {
		opts := []interface{}{&entry.TraceFunctionNameOption{*s.funcNameTracing}}
		if marshaled, err := s.logger.opts.encoderExt.EncodeExt(&e, opts...); err == nil && len(marshaled) > 0 {
			_ = s.logger.opts.writer.Write(marshaled)
		}
	} else if marshaled, err := s.logger.opts.encoder.Encode(&e); err == nil && len(marshaled) > 0 {
		_ = s.logger.opts.writer.Write(marshaled)
	}
}

func (s *_Entry) logF(level level.Level, format string, a ...interface{}) {
	e := entry.Entry{}
	e.Time = time.Now()
	e.Caller = GetCaller(true)
	if s.isRateLimited(e.Caller) {
		return
	}
	e.Level = level
	e.Msg = fmt.Sprintf(format, a...)
	e.Fields = *s.fields

	if s.ctx != nil {
		for i := range ctxFieldProviders {
			e.Fields = append(e.Fields, ctxFieldProviders[i](s.ctx)...)
		}
	}

	if s.logger.opts.encoderExt != nil {
		opts := []interface{}{&entry.TraceFunctionNameOption{*s.funcNameTracing}}

		if marshaled, err := s.logger.opts.encoderExt.EncodeExt(&e, opts...); err == nil && len(marshaled) > 0 {
			_ = s.logger.opts.writer.Write(marshaled)
		}
	} else if marshaled, err := s.logger.opts.encoder.Encode(&e); err == nil && len(marshaled) > 0 {
		_ = s.logger.opts.writer.Write(marshaled)
	}
}

func (s *_Entry) isRateLimited(caller *runtime.Frame) bool {
	if *s.interval == 0 {
		return false
	}

	key := rateLimitKey{caller.File, caller.Line}

	lastTimeMu.Lock()
	defer lastTimeMu.Unlock()
	lastTime := lastTimeMap[key]

	if time.Since(lastTime) > *s.interval {
		lastTimeMap[key] = time.Now()
		return false
	} else {
		return true
	}
}

func (s *_Entry) Trace(args ...interface{}) {
	if !s.isLevelEnabled(level.TraceLevel) {
		return
	}

	s.log(level.TraceLevel, args...)
}

func (s *_Entry) Debug(args ...interface{}) {
	if !s.isLevelEnabled(level.DebugLevel) {
		return
	}

	s.log(level.DebugLevel, args...)
}

func (s *_Entry) Info(args ...interface{}) {
	if !s.isLevelEnabled(level.InfoLevel) {
		return
	}

	s.log(level.InfoLevel, args...)
}

func (s *_Entry) Warn(args ...interface{}) {
	if !s.isLevelEnabled(level.WarnLevel) {
		return
	}

	s.log(level.WarnLevel, args...)
}

func (s *_Entry) Error(args ...interface{}) {
	if !s.isLevelEnabled(level.ErrorLevel) {
		return
	}

	s.log(level.ErrorLevel, args...)
}

func (s *_Entry) Fatal(args ...interface{}) {
	if s.isLevelEnabled(level.FatalLevel) {
		s.log(level.FatalLevel, args...)
	}

	lifecycle.Exit(1)
}

func (s *_Entry) Tracef(format string, args ...interface{}) {
	if !s.isLevelEnabled(level.TraceLevel) {
		return
	}

	s.logF(level.TraceLevel, format, args...)
}

func (s *_Entry) Debugf(format string, args ...interface{}) {
	if !s.isLevelEnabled(level.DebugLevel) {
		return
	}

	s.logF(level.DebugLevel, format, args...)
}

func (s *_Entry) Infof(format string, args ...interface{}) {
	if !s.isLevelEnabled(level.InfoLevel) {
		return
	}

	s.logF(level.InfoLevel, format, args...)
}

func (s *_Entry) Warnf(format string, args ...interface{}) {
	if !s.isLevelEnabled(level.WarnLevel) {
		return
	}

	s.logF(level.WarnLevel, format, args...)
}

func (s *_Entry) Errorf(format string, args ...interface{}) {
	if !s.isLevelEnabled(level.ErrorLevel) {
		return
	}

	s.logF(level.ErrorLevel, format, args...)
}

func (s *_Entry) Fatalf(format string, args ...interface{}) {
	if s.isLevelEnabled(level.FatalLevel) {
		s.logF(level.FatalLevel, format, args...)
	}

	lifecycle.Exit(1)
}

func (s *_Entry) Close() error {
	for _, cancel := range s.cancel {
		cancel()
	}

	return s.logger.Close()
}

func (s *_Entry) isLevelEnabled(l level.Level) bool {
	return l == level.FatalLevel || l >= level.Level(atomic.LoadUint32(s.minLevel))
}

func (s *_Entry) LogLevel() level.Level {
	return level.Level(atomic.LoadUint32(s.minLevel))
}
