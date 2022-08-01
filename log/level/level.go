package level

import (
	"fmt"
	"strings"
)

type Level uint32

const (
	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	None
)

func (level Level) String() string {
	switch level {
	case TraceLevel:
		return "TRACE"
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	case None:
		return "NONE"
	default:
		return "UNKNOWN"
	}
}

func ParseLevel(level string) (Level, error) {
	switch strings.ToLower(level) {
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "trace", "all":
		return TraceLevel, nil
	case "none":
		return None, nil
	default:
		return None, fmt.Errorf("unknown log level: %s", level)
	}
}
