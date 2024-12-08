// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/StellrisJAY/cloud-emu/platform/internal/conf"
	"github.com/StellrisJAY/cloud-emu/platform/internal/data"
	"github.com/StellrisJAY/cloud-emu/platform/internal/server"
	"github.com/StellrisJAY/cloud-emu/platform/internal/service"
	"github.com/StellrisJAY/cloud-emu/platform/internal/util"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, auth *conf.Auth, registry *conf.Registry, logger log.Logger) (*kratos.App, func(), error) {
	client := data.NewRedisClient(confData)
	dataData, cleanup, err := data.NewData(confData, client, logger)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData)
	authRepo := data.NewAuthRepo(auth)
	node := util.NewSnowflakeGenerator(confServer)
	userUseCase := biz.NewUserUseCase(userRepo, authRepo, node)
	userServer := service.NewUserService(userUseCase)
	roomRepo := data.NewRoomRepo(dataData, node)
	roomUseCase := biz.NewRoomUseCase(roomRepo, node)
	roomServer := service.NewRoomService(roomUseCase)
	roomInstanceRepo := data.NewRoomInstanceRepo(dataData)
	redsync := common.NewRedsync(client)
	iNamingClient := server.NewNacosClient(registry)
	discovery := server.NewDiscovery(iNamingClient)
	gameServerRepo := data.NewGameServerRepo(dataData, discovery)
	roomMemberRepo := data.NewRoomMemberRepo(dataData)
	roomInstanceUseCase := biz.NewRoomInstanceUseCase(roomInstanceRepo, node, redsync, gameServerRepo, roomRepo, roomMemberRepo)
	roomInstanceServer := service.NewRoomInstanceService(roomInstanceUseCase)
	grpcServer := server.NewGRPCServer(confServer, userServer, roomServer, roomInstanceServer, logger)
	httpServer := server.NewHTTPServer(confServer, auth, userServer, roomServer, roomInstanceServer, logger)
	registrar := server.NewRegistrar(iNamingClient)
	app := newApp(logger, grpcServer, httpServer, registrar)
	return app, func() {
		cleanup()
	}, nil
}
