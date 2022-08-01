package log

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/zhangel/go-frame.git/config"
	"github.com/zhangel/go-frame.git/declare"
	"github.com/zhangel/go-frame.git/lifecycle"
	"github.com/zhangel/go-frame.git/log/fields"
	"github.com/zhangel/go-frame.git/log/logger"
	"github.com/zhangel/go-frame.git/plugin"
)

const (
	loggerPrefix   = "logger"
	flagMinLevel   = "logger.min_level"
	flagRateLimit  = "logger.ratelimit"
	FlagTracingLog = "logger.rpc_tracing"
	flagPanicOnErr = "logger.panic_on_error"
)

var (
	Plugin = declare.PluginType{Name: "logger"}

	defaultLogger     logger.Logger
	ctxFieldProviders []CtxFieldsProvider
	mutex             sync.RWMutex
)

type CtxFieldsProvider func(ctx context.Context) fields.Fields

func init() {
	declare.Flags(loggerPrefix,
		declare.Flag{Name: flagMinLevel, DefaultValue: "all", Description: "Minimal log level for logger: all/trace, debug, info, warn, error, fatal, none."},
		declare.Flag{Name: flagRateLimit, DefaultValue: time.Duration(0), Description: "Logs on the same position will output only once in the interval."},
		declare.Flag{Name: FlagTracingLog, DefaultValue: false, Description: "Record rpc call tracing log."},
		declare.Flag{Name: flagPanicOnErr, DefaultValue: true, Description: "Panic while exception raised in logger handler."},
	)

	lifecycle.LifeCycle().HookFinalize(func(context.Context) {
		mutex.RLock()
		defer mutex.RUnlock()

		if defaultLogger != nil {
			_ = defaultLogger.Close()
		}
	}, lifecycle.WithName("Close default logger"), lifecycle.WithPriority(lifecycle.PriorityLowest))
}

func DefaultLogger() logger.Logger {
	if !lifecycle.LifeCycle().IsInitialized() {
		l, _ := NewConsoleLogger(true)()
		return l
	}

	mutex.RLock()
	if defaultLogger != nil {
		mutex.RUnlock()
		return defaultLogger
	}
	mutex.RUnlock()

	mutex.Lock()
	defer mutex.Unlock()

	if defaultLogger != nil {
		return defaultLogger
	}

	err := plugin.CreatePlugin(Plugin, &defaultLogger)
	if err != nil {
		log.Fatalf("[ERROR] Create logger plugin failed, err = %s.\n", err)
	}

	if defaultLogger == nil {
		defaultLogger, _ = NewConsoleLogger(true)()
	}
	defaultLogger = defaultLogger.WithConfig(config.GlobalConfig())

	return defaultLogger
}

func SetDefaultLogger(logger logger.Logger) {
	mutex.Lock()
	defer mutex.Unlock()

	defaultLogger = logger
}

func NewLoggerWithConfig(cfg config.Config) (logger.Logger, error) {
	var l logger.Logger
	err := plugin.CreatePlugin(Plugin, &l, cfg)
	if err != nil {
		return nil, err
	}

	if l == nil {
		l, _ = NewConsoleLogger(true)()
	}

	l = l.WithConfig(cfg)
	return l, nil
}

func RegisterContextFieldsProvider(provider CtxFieldsProvider) {
	mutex.Lock()
	defer mutex.Unlock()

	ctxFieldProviders = append(ctxFieldProviders, provider)
}
