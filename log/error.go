package log

import (
	"log"
	"sync"

	"github.com/zhangel/go-frame.git/config"
)

type ErrorPhase int

const (
	Init ErrorPhase = iota
	Write
)

var (
	errorHandler   func(ErrorPhase, error)
	errorHandlerMu sync.RWMutex
)

func SetLoggerErrorHandler(handler func(phase ErrorPhase, err error)) {
	errorHandlerMu.Lock()
	defer errorHandlerMu.Unlock()

	errorHandler = handler
}

func RaiseLoggerError(phase ErrorPhase, err error) {
	errorHandlerMu.RLock()
	defer errorHandlerMu.RUnlock()

	if errorHandler != nil {
		errorHandler(phase, err)
		return
	}

	switch phase {
	case Init:
		if config.Bool(flagPanicOnErr) {
			log.Fatalf("init logger failed, err = %+v", err)
		} else {
			log.Printf("init logger failed, err = %+v", err)
		}
	case Write:
	}
}
