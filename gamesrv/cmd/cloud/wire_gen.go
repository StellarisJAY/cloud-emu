// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/StellrisJAY/cloud-emu/gamesrv/internal/biz"
	"github.com/StellrisJAY/cloud-emu/gamesrv/internal/conf"
	"github.com/StellrisJAY/cloud-emu/gamesrv/internal/data"
	"github.com/StellrisJAY/cloud-emu/gamesrv/internal/server"
	"github.com/StellrisJAY/cloud-emu/gamesrv/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, registry *conf.Registry, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	gameFileRepo := data.NewGameFileRepo(dataData)
	memberAuthRepo := data.NewMemberAuthRepo(dataData)
	gameServerUseCase := biz.NewGameServerUseCase(gameFileRepo, memberAuthRepo, logger)
	gameServer := service.NewGameService(gameServerUseCase)
	grpcServer := server.NewGRPCServer(confServer, gameServer, logger)
	registrar := server.NewRegistrar(registry)
	app := newApp(logger, confServer, grpcServer, registrar)
	return app, func() {
		cleanup()
	}, nil
}
