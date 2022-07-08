// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel/sdk/trace"
	"maniverse/app/frontend/shop/internal/biz"
	"maniverse/app/frontend/shop/internal/conf"
	"maniverse/app/frontend/shop/internal/data"
	"maniverse/app/frontend/shop/internal/server"
	"maniverse/app/frontend/shop/internal/service"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(app *conf.App, confServer *conf.Server, confData *conf.Data, confService *conf.Service, auth *conf.Auth, registry *conf.Registry, logger log.Logger, tracerProvider *trace.TracerProvider) (*kratos.App, func(), error) {
	cmdable := data.NewRedisCmd(confData, logger)
	discovery := data.NewDiscovery(registry)
	userClient := data.NewUserServiceClient(discovery, auth, confService, tracerProvider)
	dataData, err := data.NewData(confData, cmdable, userClient, logger)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData, logger)
	userUseCase := biz.NewUserUseCase(userRepo, logger)
	authUseCase := biz.NewAuthUseCase(auth, userRepo)
	shopService := service.NewShopService(userUseCase, authUseCase, logger)
	httpServer := server.NewHTTPServer(confServer, auth, shopService, tracerProvider, logger)
	grpcServer := server.NewGRPCServer(confServer, auth, shopService, logger, tracerProvider)
	registrar := data.NewRegistrar(registry)
	kratosApp := newApp(logger, httpServer, grpcServer, registrar)
	return kratosApp, func() {
	}, nil
}
