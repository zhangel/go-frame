package connection

import (
	"github.com/zhangel/go-frame/hooks"

	"google.golang.org/grpc"
)

func MakeConnection(target string, unaryInt grpc.UnaryClientInterceptor, streamInt grpc.StreamClientInterceptor) (cc *grpc.ClientConn, err error) {
	return hooks.MakeConnection(target, unaryInt, streamInt)
}

func MakeConnectionWithCloseFn(target string, unaryInt grpc.UnaryClientInterceptor, streamInt grpc.StreamClientInterceptor, closeFn func()) (cc *grpc.ClientConn, err error) {
	return hooks.MakeConnectionWithCloseFn(target, unaryInt, streamInt, closeFn)
}

func IsMemoryConnection(cc *grpc.ClientConn) bool {
	return hooks.IsMemoryConnection(cc)
}
