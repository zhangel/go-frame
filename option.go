package framework

import (
	"time"

	"github.com/zhangel/go-frame/config"
)

type Options struct {
	defaultConfigSource  []config.Source
	finalizeTimeout      time.Duration
	compactUsage         bool
	flagsToShow          []string
	flagsToHide          []string
	configPreparer       func(config.Config) (config.Config, error)
	beforeConfigPreparer []func()
}

type Option func(*Options) error

func WithConfigSource(configSource ...config.Source) Option {
	return func(opts *Options) error {
		opts.defaultConfigSource = configSource
		return nil
	}
}

func WithFinalizeTimeout(timeout time.Duration) Option {
	return func(opts *Options) error {
		opts.finalizeTimeout = timeout
		return nil
	}
}

// Deprecated: please use WithCustomUsage instead.  May be removed in a future release.
func WithCompactUsage(compactUsage bool) Option {
	return func(opts *Options) error {
		opts.compactUsage = compactUsage
		return nil
	}
}

func WithCustomUsage(flagsToShow, flagsToHide []string) Option {
	return func(opts *Options) error {
		opts.flagsToShow = flagsToShow
		opts.flagsToHide = flagsToHide
		return nil
	}
}

func WithBeforeConfigPrepare(hook func()) Option {
	return func(opts *Options) error {
		opts.beforeConfigPreparer = append(opts.beforeConfigPreparer, hook)
		return nil
	}
}

func WithConfigPrepare(preparer func(config.Config) (config.Config, error)) Option {
	return func(opts *Options) error {
		opts.configPreparer = preparer
		return nil
	}
}
