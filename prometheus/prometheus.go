package prometheus

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	"google.golang.org/grpc/metadata"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"

	"github.com/modern-go/reflect2"

	"github.com/zhangel/go-frame/utils"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/zhangel/go-frame/config"
	"github.com/zhangel/go-frame/declare"
	"github.com/zhangel/go-frame/lifecycle"
	"github.com/zhangel/go-frame/log"
)

const (
	prometheusPrefix  = "Prometheus"
	flagEnable        = "enable"
	flagAddr          = "addr"
	flagPort          = "port"
	flagPath          = "path"
	flagServiceList   = "service_list"
	flagTimeHistogram = "histogram"
	flagClientMetrics = "client.metrics"

	targetMeta = "x-grpc-target"
)

var (
	prometheusConfig  config.Config
	prometheusEnabled bool
	prometheusOnce    sync.Once

	reflectMetricsType             reflect2.StructType
	reflectStartedCounter          reflect2.StructField
	reflectHandledCounter          reflect2.StructField
	reflectStreamMsgReceived       reflect2.StructField
	reflectStreamMsgSent           reflect2.StructField
	reflectHandledHistogramEnabled reflect2.StructField
	reflectHandledHistogram        reflect2.StructField

	allCodes = []codes.Code{
		codes.OK, codes.Canceled, codes.Unknown, codes.InvalidArgument, codes.DeadlineExceeded, codes.NotFound,
		codes.AlreadyExists, codes.PermissionDenied, codes.Unauthenticated, codes.ResourceExhausted,
		codes.FailedPrecondition, codes.Aborted, codes.OutOfRange, codes.Unimplemented, codes.Internal,
		codes.Unavailable, codes.DataLoss,
	}
)

func init() {
	declare.Flags(prometheusPrefix,
		declare.Flag{Name: flagEnable, DefaultValue: false, Description: "Enable prometheus metrics monitor."},
		declare.Flag{Name: flagAddr, DefaultValue: "", Description: "Address of prometheus metrics endpoint."},
		declare.Flag{Name: flagPort, DefaultValue: 9090, Description: "Port of prometheus metrics endpoint."},
		declare.Flag{Name: flagPath, DefaultValue: "/metrics", Description: "Url path of prometheus metrics endpoint."},
		declare.Flag{Name: flagServiceList, DefaultValue: "", Description: "Service list of prometheus metrics enabled, '*' means all services. (partial matching)"},
		declare.Flag{Name: flagTimeHistogram, DefaultValue: true, Description: "Enable recording of handling time of RPCs. Histogram metrics can be very expensive for Prometheus to retain and query."},
		declare.Flag{Name: flagClientMetrics, DefaultValue: false, Description: "Enable recording of client handling time of RPCs."},
	)

	lifecycle.LifeCycle().HookInitialize(func() {
		if !IsPrometheusEnabled() {
			return
		}

		metrics := &grpc_prometheus.ServerMetrics{}
		reflectMetricsType = reflect2.TypeOf(metrics).(reflect2.PtrType).Elem().(reflect2.StructType)
		reflectStartedCounter = reflectMetricsType.FieldByName("serverStartedCounter")
		reflectHandledCounter = reflectMetricsType.FieldByName("serverHandledCounter")
		reflectStreamMsgReceived = reflectMetricsType.FieldByName("serverStreamMsgReceived")
		reflectStreamMsgSent = reflectMetricsType.FieldByName("serverStreamMsgSent")
		reflectHandledHistogramEnabled = reflectMetricsType.FieldByName("serverHandledHistogramEnabled")
		reflectHandledHistogram = reflectMetricsType.FieldByName("serverHandledHistogram")

		if Config().Bool(flagTimeHistogram) {
			grpc_prometheus.EnableHandlingTimeHistogram()
			grpc_prometheus.EnableClientHandlingTimeHistogram()
		}

		startPrometheusServer()
	}, lifecycle.WithName("Enable prometheus"))
}

func startPrometheusServer() {
	go func() {
		serverAddr := Config().String(flagAddr)
		if strings.TrimSpace(serverAddr) == "" {
			serverAddr = config.String("server.addr")
			if serverAddr != "" {
				if strings.HasPrefix(serverAddr, "[::]:") {
					serverAddr = serverAddr[strings.LastIndex(serverAddr, ":"):]
				}

				if serverAddr[0] == ':' {
					serverAddr = ""
				} else {
					serverAddr, _, _ = net.SplitHostPort(serverAddr)
				}
			}

			if serverAddr == "" {
				if ip, err := utils.HostIp(); err == nil {
					serverAddr = ip + serverAddr
				}
			}
		}

		fmt.Printf("Prometheus metrics listen on http://%s:%d%s\n", serverAddr, Config().Int(flagPort), Config().String(flagPath))
		http.Handle(Config().String(flagPath), promhttp.Handler())
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", serverAddr, Config().Int(flagPort)), nil); err != nil {
			log.Fatalf("Start prometheus server failed, err = %v", err)
		}
	}()
}

func Config() config.Config {
	if prometheusConfig == nil {
		prometheusConfig = config.WithPrefix(prometheusPrefix)
	}
	return prometheusConfig
}

func IsPrometheusEnabled() bool {
	prometheusOnce.Do(func() {
		prometheusEnabled = Config().Bool(flagEnable)
	})
	return prometheusEnabled
}

func Register(server *grpc.Server) {
	// COMMENT: ??????????????????????????????????????????????????????????????????
	// 1?????????Prometheus???metrics???????????????????????????????????????????????????????????????????????????Code???Counter?????????
	// 2??????DA????????????alias???????????????????????????????????????????????????????????????????????????alias?????????????????????????????????????????????protobuf???ServiceInfo????????????????????????????????????
	// ?????????????????????????????????????????????????????????????????????metrics????????????

	//if !IsPrometheusEnabled() {
	//	return
	//}
	//
	//check := serviceListChecker()
	//serviceInfo := server.GetServiceInfo()
	//for serviceName, info := range serviceInfo {
	//	for _, mInfo := range info.Methods {
	//		if check(fmt.Sprintf("/%s/%s", serviceName, mInfo.Name)) {
	//			if err := preRegisterMethod(grpc_prometheus.DefaultServerMetrics, serviceName, &mInfo); err != nil {
	//				log.Errorf("prometheus register server failed, err = %v", err)
	//				return
	//			}
	//		}
	//	}
	//}
}

func typeFromMethodInfo(mInfo *grpc.MethodInfo) string {
	if !mInfo.IsClientStream && !mInfo.IsServerStream {
		return "unary"
	}
	if mInfo.IsClientStream && !mInfo.IsServerStream {
		return "client_stream"
	}
	if !mInfo.IsClientStream && mInfo.IsServerStream {
		return "server_stream"
	}
	return "bidi_stream"
}

func preRegisterMethod(metrics *grpc_prometheus.ServerMetrics, serviceName string, mInfo *grpc.MethodInfo) error {
	methodType := typeFromMethodInfo(mInfo)
	methodName := mInfo.Name

	if serverStartedCounter, ok := reflectStartedCounter.Get(metrics).(**prometheus.CounterVec); ok && serverStartedCounter != nil && *serverStartedCounter != nil {
		_, _ = (*serverStartedCounter).GetMetricWithLabelValues(methodType, serviceName, methodName)
	} else {
		return fmt.Errorf("register prometheus server started counter failed")
	}

	if serverStreamMsgReceived, ok := reflectStreamMsgReceived.Get(metrics).(**prometheus.CounterVec); ok && serverStreamMsgReceived != nil && *serverStreamMsgReceived != nil {
		_, _ = (*serverStreamMsgReceived).GetMetricWithLabelValues(methodType, serviceName, methodName)
	} else {
		return fmt.Errorf("register prometheus serverStreamMsgReceived failed")
	}

	if serverStreamMsgSent, ok := reflectStreamMsgSent.Get(metrics).(**prometheus.CounterVec); ok && serverStreamMsgSent != nil && *serverStreamMsgSent != nil {
		_, _ = (*serverStreamMsgSent).GetMetricWithLabelValues(methodType, serviceName, methodName)
	} else {
		return fmt.Errorf("register prometheus serverStreamMsgSent failed")
	}

	if serverHandledHistogramEnabled, ok := reflectHandledHistogramEnabled.Get(metrics).(*bool); ok && serverHandledHistogramEnabled != nil && *serverHandledHistogramEnabled {
		if serverHandledHistogram, ok := reflectHandledHistogram.Get(metrics).(**prometheus.HistogramVec); ok && serverHandledHistogram != nil && *serverHandledHistogram != nil {
			_, _ = (*serverHandledHistogram).GetMetricWithLabelValues(methodType, serviceName, methodName)
		} else {
			return fmt.Errorf("register prometheus serverHandledHistogram failed")
		}
	}

	for _, code := range allCodes {
		if serverHandledCounter, ok := reflectHandledCounter.Get(metrics).(**prometheus.CounterVec); ok && serverHandledCounter != nil && *serverHandledCounter != nil {
			_, _ = (*serverHandledCounter).GetMetricWithLabelValues(methodType, serviceName, methodName, code.String())
		} else {
			return fmt.Errorf("register prometheus serverHandledCounter failed")
		}
	}

	return nil
}

func serviceListChecker() func(string) bool {
	allowAll := false
	serviceList := Config().StringList(flagServiceList)
	for i := 0; i < len(serviceList); i++ {
		serviceList[i] = strings.ToLower(serviceList[i])
		if serviceList[i] == "*" {
			allowAll = true
			break
		}
	}

	return func(serviceName string) bool {
		if allowAll {
			return true
		}

		for _, service := range serviceList {
			if service == "" {
				continue
			}

			if strings.Contains(strings.ToLower(serviceName), service) {
				return true
			}
		}
		return false
	}
}

func replaceServiceName(ctx context.Context, fullMethod string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fullMethod
	}

	targetMd := md.Get(targetMeta)
	if len(targetMd) == 0 {
		return fullMethod
	}

	target := targetMd[len(targetMd)-1]

	path := strings.Split(fullMethod, "/")
	if len(path) != 3 {
		return fullMethod
	}

	return fmt.Sprintf("/%s/%s", target, path[2])
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	if !IsPrometheusEnabled() {
		return nil
	}

	check := serviceListChecker()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		info.FullMethod = replaceServiceName(ctx, info.FullMethod)
		if check(info.FullMethod) {
			return grpc_prometheus.UnaryServerInterceptor(ctx, req, info, handler)
		} else {
			return handler(ctx, req)
		}
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	if !IsPrometheusEnabled() {
		return nil
	}

	check := serviceListChecker()
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		info.FullMethod = replaceServiceName(ss.Context(), info.FullMethod)
		if check(info.FullMethod) {
			return grpc_prometheus.StreamServerInterceptor(srv, ss, info, handler)
		} else {
			return handler(srv, ss)
		}
	}
}

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	if !IsPrometheusEnabled() || !Config().Bool(flagClientMetrics) {
		return nil
	}

	check := serviceListChecker()
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if check(method) {
			return grpc_prometheus.UnaryClientInterceptor(ctx, method, req, reply, cc, invoker, opts...)
		} else {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
	}
}

func StreamClientInterceptor() grpc.StreamClientInterceptor {
	if !IsPrometheusEnabled() || !Config().Bool(flagClientMetrics) {
		return nil
	}

	check := serviceListChecker()
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if check(method) {
			return grpc_prometheus.StreamClientInterceptor(ctx, desc, cc, method, streamer, opts...)
		} else {
			return streamer(ctx, desc, cc, method, opts...)
		}
	}
}
