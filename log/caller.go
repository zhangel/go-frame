package log

import (
	"runtime"
	"strings"
	"sync"
)

var (
	once        sync.Once
	packageName string
	funcSkips   = map[string]interface{}{
		"Trace":     struct{}{},
		"Debug":     struct{}{},
		"Info":      struct{}{},
		"Warn":      struct{}{},
		"Error":     struct{}{},
		"Fatal":     struct{}{},
		"Tracef":    struct{}{},
		"Debugf":    struct{}{},
		"Infof":     struct{}{},
		"Warnf":     struct{}{},
		"Errorf":    struct{}{},
		"Fatalf":    struct{}{},
		"Infoln":    struct{}{},
		"Warningln": struct{}{},
		"Warningf":  struct{}{},
		"Errorln":   struct{}{},
		"Fatalln":   struct{}{},
	}
)

func GetPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

func compactFunctionName(funcName string) string {
	lastPeriod := strings.LastIndex(funcName, ".")
	if lastPeriod == -1 {
		return funcName
	} else if lastPeriod+1 >= len(funcName) {
		return funcName
	} else {
		return funcName[lastPeriod+1:]
	}
}

func GetCaller(skipLogFunc bool) *runtime.Frame {
	pcs := make([]uintptr, 10)
	depth := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	once.Do(func() {
		packageName = GetPackageName(runtime.FuncForPC(pcs[0]).Name())
	})

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := GetPackageName(f.Function)

		if pkg != packageName && (!skipLogFunc || funcSkips[compactFunctionName(f.Function)] == nil) {
			return &f
		}
	}

	return nil
}
