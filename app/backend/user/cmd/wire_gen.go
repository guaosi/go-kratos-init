// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel/sdk/trace"
	"maniverse/app/backend/user/internal/biz"
	"maniverse/app/backend/user/internal/conf"
	"maniverse/app/backend/user/internal/data"
	"maniverse/app/backend/user/internal/server"
	"maniverse/app/backend/user/internal/service"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(app *conf.App, confServer *conf.Server, confData *conf.Data, auth *conf.Auth, registry *conf.Registry, logger log.Logger, tracerProvider *trace.TracerProvider) (*kratos.App, func(), error) {
	db := data.NewDB(app, confData, logger)
	dataData, cleanup, err := data.NewData(db, logger)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData, logger)
	userUseCase := biz.NewUserUseCase(userRepo, logger)
	userService := service.NewUserService(userUseCase, logger)
	grpcServer := server.NewGRPCServer(confServer, auth, userService, logger, tracerProvider)
	registrar := server.NewRegistrar(registry)
	kratosApp := newApp(logger, grpcServer, registrar)
	return kratosApp, func() {
		cleanup()
	}, nil
}