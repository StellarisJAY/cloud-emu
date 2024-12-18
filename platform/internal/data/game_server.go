package data

import (
	"context"
	"fmt"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/go-kratos/kratos/v2/registry"
)

type GameServerRepo struct {
	data *Data
	nc   registry.Discovery
}

func NewGameServerRepo(data *Data, nc registry.Discovery) biz.GameServerRepo {
	return &GameServerRepo{
		data: data,
		nc:   nc,
	}
}

func (g *GameServerRepo) ListActiveGameServers(ctx context.Context) ([]*biz.GameServer, error) {
	serviceInstances, err := g.nc.GetService(ctx, "cloudemu-gamesrv")
	if err != nil {
		return nil, err
	}
	gameServers := make([]*biz.GameServer, len(serviceInstances))
	for i, serviceInstance := range serviceInstances {
		gameServers[i] = &biz.GameServer{
			Address: serviceInstance.Endpoints[0],
		}
	}
	return gameServers, nil
}

// OpenRoomInstance 在选中的game服务器启动房间实例
func (g *GameServerRepo) OpenRoomInstance(ctx context.Context, instance *biz.RoomInstance) (string, error) {
	client, err := common.NewGRPCClient(instance.ServerIp, int(instance.RpcPort))
	if err != nil {
		return "", err
	}
	defer client.Close()
	gameServer := v1.NewGameClient(client)
	response, err := gameServer.OpenGameInstance(ctx, &v1.OpenGameInstanceRequest{
		RoomId:     instance.RoomId,
		EmulatorId: instance.EmulatorId,
		GameId:     0,
		GameName:   "",
		GameFile:   "",
	})
	if err != nil {
		return "", err
	}
	fmt.Println(response.Data)
	return response.Data.Token, nil
}

// GetRoomInstanceToken 获取房间实例的访问token
func (g *GameServerRepo) GetRoomInstanceToken(ctx context.Context, instance *biz.RoomInstance, roomId, userId int64) (string, error) {
	client, err := common.NewGRPCClient(instance.ServerIp, int(instance.RpcPort))
	if err != nil {
		return "", err
	}
	defer client.Close()
	gameServer := v1.NewGameClient(client)
	response, err := gameServer.GetRoomInstanceAccessToken(ctx, &v1.GetRoomInstanceAccessTokenRequest{
		RoomId: roomId,
		UserId: userId,
	})
	if err != nil {
		return "", err
	}
	return response.Data.Token, nil
}
