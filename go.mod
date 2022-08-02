module github.com/zhangel/go-frame.git

go 1.16

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/ahmetb/go-linq/v3 v3.2.0
	github.com/cep21/circuit/v3 v3.2.2
	github.com/fsnotify/fsnotify v1.5.4
	github.com/go-ini/ini v1.66.6
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/iancoleman/strcase v0.2.0
	github.com/improbable-eng/grpc-web v0.15.0
	github.com/juju/ratelimit v1.0.2
	github.com/modern-go/reflect2 v1.0.2
	github.com/natefinch/npipe v0.0.0-20160621034901-c1b8fa8bdcce
	github.com/opentracing/opentracing-go v1.2.0
	github.com/orcaman/concurrent-map v1.0.0
	github.com/pelletier/go-toml v1.9.5
	github.com/prometheus/client_golang v1.12.2
	github.com/rs/cors v1.8.2
	github.com/soheilhy/cmux v0.1.5
	github.com/stretchr/testify v1.7.2
	github.com/xhit/go-str2duration/v2 v2.0.0
	go.uber.org/automaxprocs v1.5.1
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad
	google.golang.org/genproto v0.0.0-20210126160654-44e461bb6506
	google.golang.org/grpc v1.33.1
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/yaml.v3 v3.0.1
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
