package log

import (
	"github.com/zhangel/go-frame.git/declare"
	"github.com/zhangel/go-frame.git/log/encoder"
	"github.com/zhangel/go-frame.git/log/logger"
	"github.com/zhangel/go-frame.git/log/writer"
)

func init() {
	declare.Plugin(Plugin, declare.PluginInfo{Name: "stdout", Creator: NewConsoleLogger(true), ForceEnableUri: true})
	declare.Plugin(Plugin, declare.PluginInfo{Name: "stderr", Creator: NewConsoleLogger(false), ForceEnableUri: true})
}

func NewConsoleLogger(stdout bool) func() (logger.Logger, error) {
	return func() (logger.Logger, error) {
		return NewLogger(
			WithEncoder(encoder.NewSimpleTextEncoder(" ", true)),
			WithWriter(func() writer.Writer { w, _ := writer.NewConsoleWriter(stdout); return w }()),
		)
	}
}
