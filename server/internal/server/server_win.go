// +build windows

package server

import (
	"github.com/natefinch/npipe"
	"net"

	"github.com/zhangel/go-frame/server/internal/option"
)

func Listen(addr string, minPort, maxPort int) (net.Listener, error) {
	if option.IsPipeAddress(addr) {
		return npipe.Listen(addr)
	} else {
		return tcpListen(addr, minPort, maxPort)
	}
}

func defaultAddress() string {
	return "127.0.0.1"
}
