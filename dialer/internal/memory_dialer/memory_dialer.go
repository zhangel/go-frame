package memory_dialer

import (
	"context"

	"github.com/zhangel/go-frame.git/dialer/internal/option"
	"github.com/zhangel/go-frame.git/lifecycle"
	"github.com/zhangel/go-frame.git/memory_registry"
	"github.com/zhangel/go-frame.git/tracing"

	"google.golang.org/grpc"
)

var memoryDialer *memory_registry.MemoryDialer

func init() {
	lifecycle.LifeCycle().HookInitialize(func() {
		opts := []memory_registry.DialOption{}

		if tracing.DefaultTracer() != nil {
			opts = append(opts, memory_registry.DialWithUnaryInterceptor(tracing.OpenTracingClientInterceptor(tracing.DefaultTracer(), true)))
			opts = append(opts, memory_registry.DialWithStreamInterceptor(tracing.OpenTracingStreamClientInterceptor(tracing.DefaultTracer(), true)))
		}

		memoryDialer = memory_registry.NewMemoryDialer(memory_registry.GlobalRegistry, opts...)
	}, lifecycle.WithPriority(lifecycle.PriorityLowest))
}

func Dial(_ context.Context, callOpts *option.CallOptions) (*grpc.ClientConn, error) {
	if !callOpts.UseInProcDial {
		return nil, nil
	}

	return memoryDialer.Dial(callOpts.Target)
}
