// +build !windows

package option

import (
	"github.com/zhangel/go-frame/declare"
	"github.com/zhangel/go-frame/server/internal"
)

var AddrFlag = declare.Flag{Name: internal.FlagAddr, DefaultValue: internal.AutoSelectAddr, Description: "Bind address of the grpc server. optionals: address, 'auto'."}

func WithAddr(addr string) Option {
	return func(opts *Options) (err error) {
		opts.Addr = addr
		return nil
	}
}

func WithHttpAddr(addr string) Option {
	return func(opts *Options) (err error) {
		opts.HttpAddr = addr
		return nil
	}
}

func IsPipeAddress(addr string) bool {
	return false
}
