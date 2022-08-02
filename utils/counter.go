package utils

import (
	"math"
	"sync/atomic"
)

type counter struct {
	count   uint32
	current uint32
	onHit   func(...interface{})
	onMiss  func(...interface{})
}

func NewCounter(count uint32, onHit func(...interface{}), onMiss func(...interface{})) *counter {
	return &counter{
		count:   count,
		current: math.MaxUint32,
		onHit:   onHit,
		onMiss:  onMiss,
	}
}

func (s *counter) Tick(args ...interface{}) {
	if atomic.AddUint32(&s.current, 1)%s.count == 0 {
		s.onHit(args...)
	} else {
		s.onMiss(args...)
	}
}
