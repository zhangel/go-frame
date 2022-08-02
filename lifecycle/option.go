package lifecycle

import (
	"math"
	"time"
)

type Options struct {
	name     string
	priority int32
	timeout  time.Duration
}

const (
	PriorityLowest  = math.MinInt32
	PriorityLower   = math.MinInt32 >> 1
	PriorityDefault = 0
	PriorityHigher  = math.MaxInt32 >> 1
	PriorityHighest = math.MaxInt32
)

type Option func(*Options)

func generateOptions(opt ...Option) *Options {
	opts := &Options{
		priority: PriorityDefault,
		timeout:  10 * time.Second,
	}

	for _, o := range opt {
		o(opts)
	}

	return opts
}

func WithName(name string) Option {
	return func(opts *Options) {
		opts.name = name
	}
}

func WithPriority(priority int32) Option {
	return func(opts *Options) {
		opts.priority = priority
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.timeout = timeout
	}
}
