package server

import (
	"context"
	"crypto/tls"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwt2 "github.com/golang-jwt/jwt/v4"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	v1 "maniverse/api/backend/user"
	"maniverse/app/backend/user/internal/conf"
	"maniverse/app/backend/user/internal/service"
	tool_jwt "maniverse/pkg/util/jwt"
)

func NewWhiteListMatcher() selector.MatchFunc {

	whiteList := make(map[string]struct{})
	whiteList["/api.user.User/GetUserByPhone"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, ac *conf.Auth, greeter *service.UserService, logger log.Logger, tp *tracesdk.TracerProvider) *grpc.Server {
	cert, err := tls.LoadX509KeyPair(ac.TlsServerCrtPath, ac.TlsServerKeyPath)
	if err != nil {
		panic(err)
	}
	tlsConf := &tls.Config{Certificates: []tls.Certificate{cert}}
	var opts = []grpc.ServerOption{
		grpc.TLSConfig(tlsConf),
		grpc.Middleware(
			recovery.Recovery(),
			tracing.Server(
				tracing.WithTracerProvider(tp)),
			logging.Server(logger),
			selector.Server(
				jwt.Server(func(token *jwt2.Token) (interface{}, error) {
					return []byte(ac.ApiKey), nil
				}, jwt.WithSigningMethod(jwt2.SigningMethodHS256), jwt.WithClaims(func() jwt2.Claims {
					return &tool_jwt.CustomClaims{}
				})),
			).
				Match(NewWhiteListMatcher()).
				Build(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterUserServer(srv, greeter)
	return srv
}
