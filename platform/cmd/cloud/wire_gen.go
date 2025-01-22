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
func wireApp(confServer *conf.Server, confData *conf.Data, auth *conf.Auth, registry *conf.Registry, smtp *conf.Smtp, logger log.Logger) (*kratos.App, func(), error) {
	client := data.NewRedisClient(confData)
	apiClient := server.NewConsulClient(registry)
	dataData, cleanup, err := data.NewData(confData, client, logger, apiClient)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData)
	authRepo := data.NewAuthRepo(auth)
	node := util.NewSnowflakeGenerator(confServer)
	userEmailVerifyRepo := data.NewUserEmailVerifyRepo(dataData)
	emailHelper := util.NewEmailHelper(smtp)
	transaction := data.NewTransaction(dataData)
	userUseCase := biz.NewUserUseCase(userRepo, authRepo, node, userEmailVerifyRepo, emailHelper, transaction, logger)
	userServer := service.NewUserService(userUseCase)
	roomRepo := data.NewRoomRepo(dataData, node)
	roomMemberRepo := data.NewRoomMemberRepo(dataData)
	redsync := common.NewRedsync(client)
	roomInstanceRepo := data.NewRoomInstanceRepo(dataData)
	gameSaveRepo := data.NewGameSaveRepo(dataData)
	discovery := server.NewDiscovery(registry)
	gameServerRepo := data.NewGameServerRepo(dataData, discovery)
	roomUseCase := biz.NewRoomUseCase(roomRepo, node, userRepo, transaction, roomMemberRepo, redsync, logger, roomInstanceRepo, gameSaveRepo, gameServerRepo)
	notificationRepo := data.NewNotificationRepo(dataData)
	roomMemberUseCase := biz.NewRoomMemberUseCase(roomMemberRepo, roomRepo, roomInstanceRepo, notificationRepo, node, userRepo, transaction, logger)
	roomServer := service.NewRoomService(roomUseCase, roomMemberUseCase)
	emulatorRepo := data.NewEmulatorRepo(dataData, node)
	emulatorGameRepo := data.NewEmulatorGameRepo(dataData)
	roomInstanceUseCase := biz.NewRoomInstanceUseCase(roomInstanceRepo, node, redsync, gameServerRepo, roomRepo, roomMemberRepo, transaction, emulatorRepo, emulatorGameRepo, logger)
	roomInstanceServer := service.NewRoomInstanceService(roomInstanceUseCase)
	notificationUseCase := biz.NewNotificationUseCase(notificationRepo)
	notificationServer := service.NewNotificationService(notificationUseCase)
	roomMemberServer := service.NewRoomMemberService(roomMemberUseCase)
	emulatorUseCase := biz.NewEmulatorUseCase(emulatorRepo, userRepo)
	emulatorGameUseCase := biz.NewEmulatorGameUseCase(emulatorGameRepo, transaction, node, userRepo, auth, logger)
	emulatorServer := service.NewEmulatorService(emulatorUseCase, emulatorGameUseCase)
	buttonLayoutRepo := data.NewButtonLayoutRepo(dataData)
	buttonLayoutUseCase := biz.NewButtonLayoutUseCase(buttonLayoutRepo, logger)
	buttonLayoutServer := service.NewButtonLayoutService(buttonLayoutUseCase)
	keyboardBindingRepo := data.NewKeyboardBindingRepo(dataData)
	keyboardBindingUseCase := biz.NewKeyboardBindingUseCase(keyboardBindingRepo, logger)
	keyboardBindingServer := service.NewKeyboardBindingService(keyboardBindingUseCase)
	grpcServer := server.NewGRPCServer(confServer, userServer, roomServer, roomInstanceServer, notificationServer, roomMemberServer, emulatorServer, buttonLayoutServer, keyboardBindingServer, logger)
	gameSaveUseCase := biz.NewGameSaveUseCase(gameSaveRepo, roomInstanceRepo, gameServerRepo, roomMemberRepo, emulatorGameRepo, emulatorRepo, node, transaction, logger)
	gameSaveServer := service.NewGameSaveService(gameSaveUseCase)
	macroRepo := data.NewMacroRepo(dataData)
	macroUseCase := biz.NewMacroUseCase(macroRepo, emulatorRepo, roomInstanceRepo, roomMemberRepo, gameServerRepo, node, logger)
	macroServer := service.NewMacroService(macroUseCase)
	httpServer := server.NewHTTPServer(confServer, auth, userServer, roomServer, roomInstanceServer, notificationServer, roomMemberServer, emulatorServer, emulatorGameUseCase, buttonLayoutServer, keyboardBindingServer, gameSaveServer, macroServer, logger)
	registrar := server.NewRegistrar(registry)
	app := newApp(logger, confServer, grpcServer, httpServer, registrar)
	return app, func() {
		cleanup()
	}, nil
}
