package utils

import (
	"bytes"
	"runtime/pprof"
)

func FullCallStack() string {
	bytes := bytes.NewBuffer(nil)
	pprof.Lookup("goroutine").WriteTo(bytes, 1)
	return string(bytes.Bytes())
}
