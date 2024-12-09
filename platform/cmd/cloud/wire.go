//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

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
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Auth, *conf.Registry, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(util.ProviderSet, server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, common.ProviderSet, newApp))
}
