package server

import (
	"fmt"
	"github.com/StellrisJAY/cloud-emu/platform/internal/conf"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer, NewDiscovery, NewRegistrar)

func newNacosClient(c *conf.Registry) naming_client.INamingClient {
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

func newConsulClient(c *conf.Registry) *api.Client {
	client, err := api.NewClient(&api.Config{
		Address: fmt.Sprintf("%s:%d", c.ServerIp, c.Port),
	})
	if err != nil {
		panic(err)
	}
	return client
}

func NewRegistrar(c *conf.Registry) registry.Registrar {
	switch c.Scheme {
	case "nacos":
		return nacos.New(newNacosClient(c), nacos.WithPrefix(""))
	case "consul":
		return consul.New(newConsulClient(c), consul.WithHeartbeat(true), consul.WithHealthCheck(true))
	default:
		panic("unsupported scheme: " + c.Scheme)
	}
}

func NewDiscovery(c *conf.Registry) registry.Discovery {
	switch c.Scheme {
	case "nacos":
		return nacos.New(newNacosClient(c), nacos.WithPrefix(""))
	case "consul":
		return consul.New(newConsulClient(c), consul.WithHeartbeat(true), consul.WithHealthCheck(true))
	default:
		panic("unsupported scheme: " + c.Scheme)
	}
}
