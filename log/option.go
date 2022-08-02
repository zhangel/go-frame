package log

import (
	"github.com/zhangel/go-frame/log/encoder"
	"github.com/zhangel/go-frame/log/level"
	"github.com/zhangel/go-frame/log/logger"
	"github.com/zhangel/go-frame/log/writer"
)

type Options struct {
	encoder    encoder.Encoder
	encoderExt encoder.EncoderExt
	writer     writer.Writer

	minLevelOp func(logger.Logger) logger.Logger
}

type Option func(*Options) error

func WithEncoder(e encoder.Encoder) Option {
	return func(opts *Options) error {
		opts.encoder = e
		if encoderExt, ok := e.(encoder.EncoderExt); ok {
			opts.encoderExt = encoderExt
		}
		return nil
	}
}

func WithWriter(writer writer.Writer) Option {
	return func(opts *Options) error {
		opts.writer = writer
		return nil
	}
}

// Deprecated: Use Logger.WithSkipFunc instead, will deleted in future
func WithSkipLogFunc(skip bool) Option {
	return func(opts *Options) error {
		return nil
	}
}

// Deprecated: Use Logger.WithLevel instead, will deleted in future
func WithMinLevel(minLevel level.Level) Option {
	return func(opts *Options) error {
		opts.minLevelOp = func(l logger.Logger) logger.Logger {
			return l.WithLevel(minLevel)
		}
		return nil
	}
}
