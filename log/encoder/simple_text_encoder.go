package encoder

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/zhangel/go-frame/log/entry"
)

type SimpleTextEncoder struct {
	sep            string
	withFixedParts bool
}

func NewSimpleTextEncoder(sep string, withFixedParts bool) Encoder {
	return &SimpleTextEncoder{sep, withFixedParts}
}

func (s *SimpleTextEncoder) Encode(entry *entry.Entry) ([]byte, error) {
	return s.EncodeExt(entry)
}

func (s *SimpleTextEncoder) EncodeExt(e *entry.Entry, opts ...interface{}) ([]byte, error) {
	var logItem []string

	enableFunctionNameTracing := false
	for _, opt := range opts {
		switch o := opt.(type) {
		case *entry.TraceFunctionNameOption:
			enableFunctionNameTracing = o.Enabled
		}
	}

	if s.withFixedParts {
		_, fileName := filepath.Split(e.Caller.File)
		logItem = append(logItem, fmt.Sprintf("%s", e.Time.Format("2006-01-02 15:04:05.000000")))
		logItem = append(logItem, fmt.Sprintf("[%s]", e.Level.String()))
		if enableFunctionNameTracing {
			_, funcName := filepath.Split(e.Caller.Function)
			logItem = append(logItem, fmt.Sprintf("<%s:%d:%s>", fileName, e.Caller.Line, funcName))
		} else {
			logItem = append(logItem, fmt.Sprintf("<%s:%d>", fileName, e.Caller.Line))
		}
	}

	for _, field := range e.Fields {
		logItem = append(logItem, fmt.Sprintf("%s:%v", field.K, field.V))
	}

	logItem = append(logItem, fmt.Sprintf("msg:%s", e.Msg))

	return []byte(strings.Join(logItem, s.sep)), nil
}
