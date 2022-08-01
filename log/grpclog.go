package log

import (
	"github.com/zhangel/go-frame.git/lifecycle"
	"google.golang.org/grpc/grpclog"
)

type grpcLogger struct{}

func init() {
	lifecycle.LifeCycle().HookInitialize(func() {
		grpclog.SetLoggerV2(&grpcLogger{})
	}, lifecycle.WithName("Set grpc logger"))
}

func (s *grpcLogger) Info(args ...interface{}) {
	Debug(args...)
}

func (s *grpcLogger) Infoln(args ...interface{}) {
	Debug(args...)
}

func (s *grpcLogger) Infof(format string, args ...interface{}) {
	Debugf(format, args...)
}

func (s *grpcLogger) Warning(args ...interface{}) {
	Warn(args...)
}

func (s *grpcLogger) Warningln(args ...interface{}) {
	Warn(args...)
}

func (s *grpcLogger) Warningf(format string, args ...interface{}) {
	Warnf(format, args...)
}

func (s *grpcLogger) Error(args ...interface{}) {
	Error(args...)
}

func (s *grpcLogger) Errorln(args ...interface{}) {
	Error(args...)
}

func (s *grpcLogger) Errorf(format string, args ...interface{}) {
	Errorf(format, args...)
}

func (s *grpcLogger) Fatal(args ...interface{}) {
	Fatal(args...)
}

func (s *grpcLogger) Fatalln(args ...interface{}) {
	Fatal(args...)
}

func (s *grpcLogger) Fatalf(format string, args ...interface{}) {
	Fatalf(format, args...)
}

func (s *grpcLogger) V(l int) bool {
	return true
}
