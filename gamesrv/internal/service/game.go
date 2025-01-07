package service

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/gamesrv/internal/biz"
)

type GameService struct {
	v1.UnimplementedGameServer
	uc *biz.GameServerUseCase
}

func NewGameService(uc *biz.GameServerUseCase) v1.GameServer {
	return &GameService{
		uc: uc,
	}
}

func (g *GameService) OpenGameInstance(ctx context.Context, request *v1.OpenGameInstanceRequest) (*v1.OpenGameInstanceResponse, error) {
	params := biz.CreateRoomInstanceParams{
		RoomId: request.RoomId,
		Auth: &biz.MemberAuthInfo{
			UserId: request.Auth.UserId,
			AppId:  request.Auth.AppId,
			Ip:     request.Auth.Ip,
		},
		EmulatorId:   request.EmulatorId,
		GameId:       request.GameId,
		GameData:     request.GameData,
		EmulatorType: request.EmulatorType,
	}
	token, sessionKey, err := g.uc.CreateRoomInstance(ctx, params)
	if err != nil {
		return nil, err
	}
	return &v1.OpenGameInstanceResponse{
		Code:    200,
		Message: "创建成功",
		Data:    &v1.GetGameInstanceResult{Token: token, SessionKey: sessionKey},
	}, nil
}

func (g *GameService) GetRoomInstanceAccessToken(ctx context.Context, request *v1.GetRoomInstanceAccessTokenRequest) (*v1.GetRoomInstanceAccessTokenResponse, error) {
	token, err := g.uc.GetRoomInstanceToken(ctx, request.RoomId, &biz.MemberAuthInfo{
		UserId: request.Auth.UserId,
		AppId:  request.Auth.AppId,
		Ip:     request.Auth.Ip,
	})
	if err != nil {
		return nil, err
	}
	return &v1.GetRoomInstanceAccessTokenResponse{
		Code:    200,
		Message: "创建成功",
		Data:    &v1.GetGameInstanceResult{Token: token},
	}, nil
}

func (g *GameService) ShutdownRoomInstance(ctx context.Context, request *v1.ShutdownRoomInstanceRequest) (*v1.ShutdownRoomInstanceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GameService) OpenGameConnection(ctx context.Context, request *v1.GameSrvOpenGameConnectionRequest) (*v1.GameSrvOpenGameConnectionResponse, error) {
	sdpOffer, err := g.uc.OpenGameConnection(ctx, request.RoomId, request.Token, &biz.MemberAuthInfo{
		UserId: request.Auth.UserId,
		AppId:  request.Auth.AppId,
		Ip:     request.Auth.Ip,
	})
	if err != nil {
		return nil, err
	}
	return &v1.GameSrvOpenGameConnectionResponse{
		Code:    200,
		Message: "创建成功",
		Data:    &v1.GameSrvOpenGameConnectionResult{SdpOffer: sdpOffer},
	}, nil
}

func (g *GameService) SdpAnswer(ctx context.Context, request *v1.GameSrvSdpAnswerRequest) (*v1.GameSrvSdpAnswerResponse, error) {
	err := g.uc.SdpAnswer(ctx, request.RoomId, request.Token, &biz.MemberAuthInfo{
		UserId: request.Auth.UserId,
		AppId:  request.Auth.AppId,
		Ip:     request.Auth.Ip,
	}, request.SdpAnswer)
	if err != nil {
		return nil, err
	}
	return &v1.GameSrvSdpAnswerResponse{
		Code:    200,
		Message: "创建成功",
	}, nil
}

func (g *GameService) AddIceCandidate(ctx context.Context, request *v1.GameSrvAddIceCandidateRequest) (*v1.GameSrvAddIceCandidateResponse, error) {
	err := g.uc.AddICECandidate(ctx, request.RoomId, request.Token, &biz.MemberAuthInfo{
		UserId: request.Auth.UserId,
		AppId:  request.Auth.AppId,
		Ip:     request.Auth.Ip,
	}, request.IceCandidate)
	if err != nil {
		return nil, err
	}
	return &v1.GameSrvAddIceCandidateResponse{
		Code:    200,
		Message: "创建成功",
	}, nil
}

func (g *GameService) GetIceCandidate(ctx context.Context, request *v1.GetIceCandidateRequest) (*v1.GetIceCandidateResponse, error) {
	candidates, err := g.uc.GetLocalICECandidate(ctx, request.RoomId, request.Token, &biz.MemberAuthInfo{
		UserId: request.Auth.UserId,
		AppId:  request.Auth.AppId,
		Ip:     request.Auth.Ip,
	})
	if err != nil {
		return nil, err
	}
	return &v1.GetIceCandidateResponse{
		Code:       200,
		Message:    "创建成功",
		Candidates: candidates,
	}, nil
}

func (g *GameService) RestartGameInstance(ctx context.Context, request *v1.RestartGameInstanceRequest) (*v1.RestartGameInstanceResponse, error) {
	err := g.uc.Restart(ctx, biz.RestartParams{
		RoomId:       request.RoomId,
		UserId:       request.UserId,
		EmulatorType: request.EmulatorType,
		GameName:     request.GameName,
		GameUrl:      request.GameUrl,
		EmulatorId:   request.EmulatorId,
		GameId:       request.GameId,
		GameData:     request.GameData,
	})
	if err != nil {
		return nil, err
	}
	return &v1.RestartGameInstanceResponse{Code: 200, Message: "重启成功"}, nil
}

func (g *GameService) SaveGame(ctx context.Context, request *v1.GameSrvSaveGameRequest) (*v1.GameSrvSaveGameResponse, error) {
	emulatorId, gameId, data, err := g.uc.SaveGame(ctx, request.RoomId, request.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.GameSrvSaveGameResponse{
		Code:    200,
		Message: "保存成功",
		Data: &v1.GameSrvSaveResult{
			EmulatorId: emulatorId,
			GameId:     gameId,
			SaveData:   data,
		},
	}, nil
}

func (g *GameService) LoadSave(ctx context.Context, request *v1.GameSrvLoadSaveRequest) (*v1.GameSrvLoadSaveResponse, error) {
	err := g.uc.LoadSave(ctx, biz.LoadSaveParams{
		RoomId:       request.RoomId,
		UserId:       request.UserId,
		EmulatorType: request.EmulatorType,
		GameName:     request.GameName,
		EmulatorId:   request.EmulatorId,
		GameId:       request.GameId,
		GameData:     request.GameData,
		SaveData:     request.SaveData,
	})

	if err != nil {
		return nil, err
	}
	return &v1.GameSrvLoadSaveResponse{
		Code:    200,
		Message: "加载成功",
	}, nil
}

func (g *GameService) ListOnlineRoomMember(ctx context.Context, request *v1.ListOnlineRoomMemberRequest) (*v1.ListOnlineRoomMemberResponse, error) {
	memberIdList, err := g.uc.ListOnlineRoomMember(ctx, request.RoomId)
	if err != nil {
		return nil, err
	}
	return &v1.ListOnlineRoomMemberResponse{
		Code:          200,
		Message:       "获取成功",
		RoomMemberIds: memberIdList,
	}, nil
}

func (g *GameService) GetControllerPlayers(ctx context.Context, request *v1.GameSrvGetControllerPlayersRequest) (*v1.GameSrvGetControllerPlayersResponse, error) {
	cps, err := g.uc.GetControllerPlayers(ctx, request.RoomId)
	if err != nil {
		return nil, err
	}
	result := make([]*v1.GameSrvControllerPlayer, len(cps))
	for i, cp := range cps {
		result[i] = &v1.GameSrvControllerPlayer{
			UserId:       cp.UserId,
			ControllerId: int32(cp.ControllerId),
			Label:        cp.Label,
		}
	}
	return &v1.GameSrvGetControllerPlayersResponse{
		Code:    200,
		Message: "获取成功",
		Data:    result,
	}, nil
}
