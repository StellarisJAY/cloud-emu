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
	return &v1.OpenGameInstanceResponse{
		Code:    200,
		Message: "创建成功",
		Data:    &v1.GetGameInstanceResult{Token: "a"},
	}, nil
}

func (g *GameService) GetRoomInstanceAccessToken(ctx context.Context, request *v1.GetRoomInstanceAccessTokenRequest) (*v1.GetRoomInstanceAccessTokenResponse, error) {
	return &v1.GetRoomInstanceAccessTokenResponse{
		Code:    200,
		Message: "创建成功",
		Data:    &v1.GetGameInstanceResult{Token: "a"},
	}, nil
}
