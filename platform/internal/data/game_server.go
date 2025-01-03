package data

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
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
func (g *GameServerRepo) OpenRoomInstance(ctx context.Context, instance *biz.RoomInstance, auth biz.RoomMemberAuth, gameData []byte) (string, string, error) {
	client, err := common.NewGRPCClient(instance.ServerIp, int(instance.RpcPort))
	if err != nil {
		return "", "", err
	}
	defer client.Close()
	gameServer := v1.NewGameClient(client)
	response, err := gameServer.OpenGameInstance(ctx, &v1.OpenGameInstanceRequest{
		RoomId: instance.RoomId,
		Auth: &v1.RoomMemberAuth{
			UserId: auth.UserId,
			Ip:     auth.Ip,
			AppId:  auth.AppId,
		},
		GameData:     gameData,
		EmulatorType: instance.EmulatorType,
		EmulatorId:   instance.EmulatorId,
		GameId:       instance.GameId,
	})
	if err != nil {
		return "", "", err
	}
	if response.Code != 200 {
		return "", "", errors.New(int(response.Code), "Service Error", response.Message)
	}
	return response.Data.Token, response.Data.SessionKey, nil
}

// GetRoomInstanceToken 获取房间实例的访问token
func (g *GameServerRepo) GetRoomInstanceToken(ctx context.Context, instance *biz.RoomInstance, roomId int64, auth biz.RoomMemberAuth) (string, error) {
	client, err := common.NewGRPCClient(instance.ServerIp, int(instance.RpcPort))
	if err != nil {
		return "", err
	}
	defer client.Close()
	gameServer := v1.NewGameClient(client)
	response, err := gameServer.GetRoomInstanceAccessToken(ctx, &v1.GetRoomInstanceAccessTokenRequest{
		RoomId: roomId,
		Auth: &v1.RoomMemberAuth{
			UserId: auth.UserId,
			Ip:     auth.Ip,
			AppId:  auth.AppId,
		},
	})
	if err != nil {
		return "", err
	}
	if response.Code != 200 {
		return "", errors.New(int(response.Code), "Service Error", response.Message)
	}
	return response.Data.Token, nil
}

func (g *GameServerRepo) OpenGameConnection(ctx context.Context, instance *biz.RoomInstance, token string, auth biz.RoomMemberAuth) (string, error) {
	client, err := common.NewGRPCClient(instance.ServerIp, int(instance.RpcPort))
	if err != nil {
		return "", err
	}
	defer client.Close()
	gameServer := v1.NewGameClient(client)
	resp, err := gameServer.OpenGameConnection(ctx, &v1.GameSrvOpenGameConnectionRequest{
		RoomId: instance.RoomId,
		Token:  token,
		Auth: &v1.RoomMemberAuth{
			UserId: auth.UserId,
			Ip:     auth.Ip,
			AppId:  auth.AppId,
		},
	})
	if err != nil {
		return "", err
	}
	if resp.Code != 200 {
		return "", errors.New(int(resp.Code), "Service Error", resp.Message)
	}
	return resp.Data.SdpOffer, nil
}

func (g *GameServerRepo) SdpAnswer(ctx context.Context, instance *biz.RoomInstance, token string, auth biz.RoomMemberAuth, sdpAnswer string) error {
	client, err := common.NewGRPCClient(instance.ServerIp, int(instance.RpcPort))
	if err != nil {
		return err
	}
	defer client.Close()
	gameServer := v1.NewGameClient(client)
	resp, err := gameServer.SdpAnswer(ctx, &v1.GameSrvSdpAnswerRequest{
		RoomId: instance.RoomId,
		Token:  token,
		Auth: &v1.RoomMemberAuth{
			UserId: auth.UserId,
			Ip:     auth.Ip,
			AppId:  auth.AppId,
		},
		SdpAnswer: sdpAnswer,
	})
	if err != nil {
		return err
	}
	if resp.Code != 200 {
		return errors.New(int(resp.Code), "Service Error", resp.Message)
	}
	return nil
}

func (g *GameServerRepo) AddICECandidate(ctx context.Context, instance *biz.RoomInstance, token string, auth biz.RoomMemberAuth, candidate string) error {
	client, err := common.NewGRPCClient(instance.ServerIp, int(instance.RpcPort))
	if err != nil {
		return err
	}
	defer client.Close()
	gameServer := v1.NewGameClient(client)
	resp, err := gameServer.AddIceCandidate(ctx, &v1.GameSrvAddIceCandidateRequest{
		RoomId: instance.RoomId,
		Token:  token,
		Auth: &v1.RoomMemberAuth{
			UserId: auth.UserId,
			Ip:     auth.Ip,
			AppId:  auth.AppId,
		},
		IceCandidate: candidate,
	})
	if err != nil {
		return err
	}
	if resp.Code != 200 {
		return errors.New(int(resp.Code), "Service Error", resp.Message)
	}
	return nil
}

func (g *GameServerRepo) GetServerICECandidate(ctx context.Context, instance *biz.RoomInstance, token string, auth biz.RoomMemberAuth) ([]string, error) {
	client, err := common.NewGRPCClient(instance.ServerIp, int(instance.RpcPort))
	if err != nil {
		return nil, err
	}
	defer client.Close()
	gameServer := v1.NewGameClient(client)
	resp, err := gameServer.GetIceCandidate(ctx, &v1.GetIceCandidateRequest{
		RoomId: instance.RoomId,
		Token:  token,
		Auth: &v1.RoomMemberAuth{
			UserId: auth.UserId,
			Ip:     auth.Ip,
			AppId:  auth.AppId,
		},
	})
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, errors.New(int(resp.Code), "Service Error", resp.Message)
	}
	return resp.Candidates, nil
}

func (g *GameServerRepo) RestartGameInstance(ctx context.Context, instance *biz.RoomInstance, params biz.RestartParams) error {
	client, err := common.NewGRPCClient(instance.ServerIp, int(instance.RpcPort))
	if err != nil {
		return err
	}
	defer client.Close()
	gameServer := v1.NewGameClient(client)
	resp, err := gameServer.RestartGameInstance(ctx, &v1.RestartGameInstanceRequest{
		RoomId:       instance.RoomId,
		UserId:       params.UserId,
		EmulatorType: params.EmulatorType,
		GameName:     params.GameName,
		EmulatorId:   params.EmulatorId,
		GameId:       params.GameId,
		GameData:     params.GameData,
	})

	if err != nil {
		return err
	}
	if resp.Code != 200 {
		return errors.New(int(resp.Code), "Service Error", resp.Message)
	}
	return nil
}

func (g *GameServerRepo) SaveGame(ctx context.Context, instance *biz.RoomInstance, roomId, userId int64) (int64, int64, []byte, error) {
	client, err := common.NewGRPCClient(instance.ServerIp, int(instance.RpcPort))
	if err != nil {
		return 0, 0, nil, err
	}
	defer client.Close()
	gameServer := v1.NewGameClient(client)
	resp, err := gameServer.SaveGame(ctx, &v1.GameSrvSaveGameRequest{
		RoomId: roomId,
		UserId: userId,
	})
	if err != nil {
		return 0, 0, nil, err
	}
	if resp.Code != 200 {
		return 0, 0, nil, errors.New(int(resp.Code), "Service Error", resp.Message)
	}
	return resp.Data.EmulatorId, resp.Data.GameId, resp.Data.SaveData, nil
}

func (g *GameServerRepo) LoadSave(ctx context.Context, instance *biz.RoomInstance, params biz.LoadSaveParams) error {
	client, err := common.NewGRPCClient(instance.ServerIp, int(instance.RpcPort))
	if err != nil {
		return err
	}
	defer client.Close()
	gameServer := v1.NewGameClient(client)
	resp, err := gameServer.LoadSave(ctx, &v1.GameSrvLoadSaveRequest{
		RoomId:       instance.RoomId,
		UserId:       params.UserId,
		EmulatorId:   params.EmulatorId,
		EmulatorType: params.EmulatorType,
		GameId:       params.GameId,
		GameName:     params.GameName,
		GameData:     params.GameData,
		SaveData:     params.SaveData,
	})
	if err != nil {
		return err
	}
	if resp.Code != 200 {
		return errors.New(int(resp.Code), "Service Error", resp.Message)
	}
	return nil
}
