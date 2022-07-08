package data

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-redis/redis/v8"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"maniverse/api/backend/user"
	"maniverse/app/frontend/shop/internal/conf"
	"os"
	"time"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewRedisCmd, NewDiscovery, NewRegistrar, NewUserServiceClient, NewUserRepo)

// Data .
type Data struct {
	log      *log.Helper
	uc       user.UserClient
	redisCli redis.Cmdable
	// TODO wrapped database client
}

// NewData .
func NewData(c *conf.Data, redisCmd redis.Cmdable, uc user.UserClient, logger log.Logger) (*Data, error) {
	log := log.NewHelper(log.With(logger, "module", "shop-service/data"))

	d := &Data{
		log:      log,
		uc:       uc,
		redisCli: redisCmd,
	}

	return d, nil
}

func NewRedisCmd(conf *conf.Data, logger log.Logger) redis.Cmdable {
	log := log.NewHelper(log.With(logger, "module", "shop-service/data/redis"))
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
		DialTimeout:  time.Second * 2,
		PoolSize:     10,
	})
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	err := client.Ping(timeout).Err()
	if err != nil {
		log.Fatalf("redis connect error: %v", err)
	}
	return client
}

func NewDiscovery(conf *conf.Registry) registry.Discovery {
	c := consulAPI.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r
}

func NewRegistrar(conf *conf.Registry) registry.Registrar {
	c := consulAPI.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r
}
func NewUserServiceClient(r registry.Discovery, confAuth *conf.Auth, confService *conf.Service, tp *tracesdk.TracerProvider) user.UserClient {
	// Load CA certificate pem file.
	b, err := os.ReadFile(confAuth.TlsCaCrtPath)
	if err != nil {
		panic(err)
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		panic(err)
	}
	tlsConf := &tls.Config{ServerName: confAuth.TlsServerName, RootCAs: cp}
	conn, err := grpc.Dial(
		context.Background(),
		grpc.WithTLSConfig(tlsConf),
		grpc.WithEndpoint(fmt.Sprintf("discovery:///%s", confService.User)),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			tracing.Client(tracing.WithTracerProvider(tp)),
			recovery.Recovery(),
			jwt.Client(func(token *jwt2.Token) (interface{}, error) {
				return []byte(confAuth.UserServiceKey), nil
			}, jwt.WithSigningMethod(jwt2.SigningMethodHS256)),
		),
	)
	if err != nil {
		panic(err)
	}
	c := user.NewUserClient(conn)
	return c
}
