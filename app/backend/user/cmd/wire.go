// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"maniverse/app/backend/user/internal/biz"
	"maniverse/app/backend/user/internal/conf"
	"maniverse/app/backend/user/internal/data"
	"maniverse/app/backend/user/internal/server"
	"maniverse/app/backend/user/internal/service"
)

// wireApp init kratos application.
func wireApp(*conf.App, *conf.Server, *conf.Data, *conf.Auth, *conf.Registry, log.Logger, *tracesdk.TracerProvider) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
