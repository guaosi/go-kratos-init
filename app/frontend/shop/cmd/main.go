package main

import (
	"flag"
	"github.com/go-kratos/kratos/v2/registry"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.4.0"
	"maniverse/pkg/util/consul"
	tool_file "maniverse/pkg/util/file"
	"maniverse/pkg/util/zap"
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"maniverse/app/frontend/shop/internal/conf"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "frontend.guaosi.shop.service"
	// Version is the version of the compiled software.
	Version = "v1"
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, hs *http.Server, gs *grpc.Server, rr registry.Registrar) *kratos.App {
	return kratos.New(
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
		kratos.Registrar(rr),
	)
}

func main() {
	flag.Parse()
	// 配置文件初始化
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()
	// 读取本地配置文件，获取配置中心地址
	if err := c.Load(); err != nil {
		panic(err)
	}

	var rc conf.Registry
	if err := c.Scan(&rc); err != nil {
		panic(err)
	}

	// 解析配置中心中的配置信息
	consulClient, err := consul.GetConsulConfiguration(rc.Consul.Address, rc.Consul.Path)
	if err != nil {
		panic(err)
	}
	var bc conf.Bootstrap
	if err = consulClient.Scan(&bc); err != nil {
		panic(err)
	}
	if bc.App.Name != "" {
		Name = bc.App.Name
	}
	if bc.App.Version != "" {
		Version = bc.App.Version
	}
	if bc.App.LogPath != "" {
		if !strings.HasSuffix(bc.App.LogPath, "/") {
			bc.App.LogPath = bc.App.LogPath + "/"
		}
		if err = tool_file.CreateIfNotExist(bc.App.LogPath); err != nil {
			panic(err)
		}
	}
	// 日志初始化
	logger := log.With(zap.NewLogger(Name, bc.App.LogPath),
		"ts", log.Timestamp("2006-01-02 15:04:05.999"),
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	// 链路追踪初始化
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(bc.Trace.Endpoint)))
	if err != nil {
		panic(err)
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(Name),
		)),
	)
	// 开启主程序
	app, cleanup, err := wireApp(bc.App, bc.Server, bc.Data, bc.Service, bc.Auth, &rc, logger, tp)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
