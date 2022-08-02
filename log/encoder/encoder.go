package encoder

import (
	"github.com/zhangel/go-frame/log/entry"
)

type Encoder interface {
	Encode(entry *entry.Entry) ([]byte, error)
}

type EncoderExt interface {
	EncodeExt(entry *entry.Entry, opts ...interface{}) ([]byte, error)
}
