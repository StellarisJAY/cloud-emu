package server

import (
	"github.com/StellrisJAY/cloud-emu/gamesrv/internal/conf"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewNacosClient, NewRegistrar, NewDiscovery)

func NewNacosClient(c *conf.Registry) naming_client.INamingClient {
	cliConfig := constant.ClientConfig{
		Username:    c.UserName,
		Password:    c.Password,
		NamespaceId: c.NamespaceId,
	}
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig: &cliConfig,
			ServerConfigs: []constant.ServerConfig{
				{
					Scheme: "http",
					IpAddr: c.ServerIp,
					Port:   uint64(c.Port),
				},
			},
		},
	)
	if err != nil {
		panic(err)
	}
	return client
}

func NewRegistrar(cli naming_client.INamingClient) registry.Registrar {
	return nacos.New(cli, nacos.WithPrefix(""))
}

func NewDiscovery(cli naming_client.INamingClient) registry.Discovery {
	return nacos.New(cli, nacos.WithPrefix(""))
}
