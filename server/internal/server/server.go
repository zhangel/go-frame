package server

import (
	"context"
	"fmt"
	"github.com/zhangel/go-frame/lifecycle"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"

	"github.com/zhangel/go-frame/log"
	"github.com/zhangel/go-frame/memory_registry"
	"github.com/zhangel/go-frame/server/internal"
	"github.com/zhangel/go-frame/server/internal/grpc_server"
	"github.com/zhangel/go-frame/server/internal/http_server"
	"github.com/zhangel/go-frame/server/internal/option"
	"github.com/zhangel/go-frame/server/internal/service"
	"github.com/zhangel/go-frame/utils"
)

type Server struct {
	grpcServiceDesc grpc_server.ServiceDescList
	httpServiceDesc http_server.ServiceDescList
	interrupting    uint32
}

func (s *Server) RegisterService(srvRegister interface{}, srvHandler service.Provider, opt ...option.RegisterOption) error {
	return s.grpcServiceDesc.RegisterServiceHandler(srvRegister, srvHandler, opt...)
}

func (s *Server) RegisterHttpService(route http_server.Route, handler http.Handler, opt ...option.HttpRegisterOption) error {
	return s.httpServiceDesc.RegisterServiceHandler(route, handler, opt...)
}

func (s *Server) Run(opts *option.Options) (err error) {
	defer func() {
		if atomic.LoadUint32(&s.interrupting) == 1 {
			err = nil
		}
	}()

	if opts.OptionsHook != nil {
		opts.OptionsHook(opts)
	}

	lifecycle.LifeCycle().HookFinalize(func(context.Context) {
		atomic.StoreUint32(&s.interrupting, 1)
	}, lifecycle.WithName("Set server interrupt flag"), lifecycle.WithPriority(lifecycle.PriorityHighest))

	return s.serve(opts)
}

func (s *Server) serve(opts *option.Options) error {
	srv := grpc.NewServer(grpc_server.PrepareGrpcOptions(opts)...)
	grpcServer := grpc_server.NewGrpcServer(srv, opts, s.grpcServiceDesc)

	var httpServer *http_server.Server
	if opts.EnableHttp || opts.EnableGrpcWeb {
		httpServer = http_server.NewHttpServer(srv, opts, s.grpcServiceDesc, s.httpServiceDesc)
	}

	listener, err := Listen(opts.Addr, opts.MinPort, opts.MaxPort)
	if err != nil {
		log.Fatal("[ERROR]", err)
		return nil
	}

	if httpServer == nil {
		return grpcServer.Serve(listener, grpcServer.Prepare)
	}

	if opts.MultiplexAddr {
		if opts.WithTls {
			// ??????????????????TLS???????????????????????????cmux??????????????????cmux????????????ALPN????????????????????????http2?????????????????????HttpServer????????????HttpServer?????????Grpc??????????????????
			// ?????????TLS????????????????????????HttpServer??????GRPC???????????????????????????????????????TLS?????????????????????????????????????????????????????????
			return httpServer.Serve(listener, func(l net.Listener) func() {
				grpcFinalizer := grpcServer.Prepare(l)
				httpFinalizer := httpServer.Prepare(l)
				return func() {
					var wg sync.WaitGroup
					wg.Add(2)
					go func() {
						defer wg.Done()
						grpcFinalizer()
					}()
					go func() {
						defer wg.Done()
						httpFinalizer()
					}()
					wg.Wait()
				}
			})
		} else {
			m := cmux.New(listener)
			m.SetReadTimeout(1 * time.Second)

			httpListener := m.Match(cmux.HTTP1Fast("PATCH", "LINK", "UNLINK", "PURGE", "VIEW", "PROPFIND"))
			grpcListener := m.Match(cmux.Any())
			go func() {
				if err := grpcServer.Serve(grpcListener, grpcServer.Prepare); err != nil {
					if atomic.LoadUint32(&s.interrupting) != 1 {
						log.Fatal("[ERROR]", err)
					}
				}
			}()

			go func() {
				if err := httpServer.Serve(httpListener, httpServer.Prepare); err != nil {
					if atomic.LoadUint32(&s.interrupting) != 1 {
						log.Fatal("[ERROR]", err)
					}
				}
			}()

			return m.Serve()
		}
	} else {
		httpListener, err := Listen(opts.HttpAddr, opts.MinPort, opts.MaxPort)
		if err != nil {
			log.Fatal("[ERROR]", err)
			return nil
		}

		go func() {
			if err := httpServer.Serve(httpListener, httpServer.Prepare); err != nil {
				if atomic.LoadUint32(&s.interrupting) != 1 {
					log.Fatal("[ERROR]", err)
				}
			}
		}()

		return grpcServer.Serve(listener, grpcServer.Prepare)
	}
}

func tcpListen(addr string, minPort, maxPort int) (net.Listener, error) {
	if strings.ToLower(strings.TrimSpace(addr)) == internal.AutoSelectAddr {
		addr = defaultAddress() + ":0"
	}

	if strings.HasPrefix(addr, "[::]:") {
		addr = addr[strings.LastIndex(addr, ":"):]
	} else if strings.HasPrefix(addr, "0.0.0.0") {
		if idx := strings.LastIndex(addr, ":"); idx != -1 {
			addr = addr[idx:]
		} else {
			addr = ""
		}
	}

	if addr == "" {
		if ip, err := utils.HostIp(); err != nil {
			return nil, fmt.Errorf("invalid address specified, err = %v", err)
		} else {
			addr = ip + ":0"
		}
	} else if addr[0] == ':' {
		if ip, err := utils.HostIp(); err != nil {
			return nil, fmt.Errorf("invalid address specified, err = %v", err)
		} else {
			addr = ip + addr
		}
	} else if !strings.Contains(addr, ":") {
		addr = addr + ":0"
	}

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid address specified, addr = %s, err = %v", addr, err)
	}

	if ip := net.ParseIP(host); ip == nil {
		return nil, fmt.Errorf("invalid address specified, addr = %s, err = %v", addr, err)
	}

	if option.IsAutoPort(addr) {
		if minPort == -1 {
			minPort = 10000
		}
		if maxPort == -1 {
			maxPort = 20000
		}

		if maxPort < minPort {
			return nil, fmt.Errorf("invalid port range, max port is less than min port")
		}

		rand.Seed(time.Now().UnixNano())
		for i := 0; i < 50; i++ {
			port := int(int32(minPort) + rand.Int31n(int32(maxPort)-int32(minPort)+1))
			l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
			if err == nil {
				return l, nil
			}
		}

		return nil, fmt.Errorf("no available port found in port range [%d - %d]", minPort, maxPort)
	}

	if l, err := net.Listen("tcp", addr); err != nil {
		return nil, err
	} else {
		return l, nil
	}
}

func (s *Server) ServiceInfos() ([]*memory_registry.ServiceInfo, error) {
	var result []*memory_registry.ServiceInfo

	for _, sd := range s.grpcServiceDesc {
		if serviceInfo, err := memory_registry.GetServiceInfo(sd.Register, sd.Handler); err != nil {
			return nil, err
		} else {
			result = append(result, serviceInfo)
		}
	}

	return result, nil
}
