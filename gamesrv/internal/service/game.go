package service

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/api/v1"
)

type GameService struct {
	v1.UnimplementedGameServer
}

func NewGameService() v1.GameServer {
	return &GameService{}
}

func (g *GameService) OpenGameInstance(ctx context.Context, request *v1.OpenGameInstanceRequest) (*v1.OpenGameInstanceResponse, error) {
	return &v1.OpenGameInstanceResponse{}, nil
}

func (g *GameService) GetRoomInstanceAccessToken(ctx context.Context, request *v1.GetRoomInstanceAccessTokenRequest) (*v1.GetRoomInstanceAccessTokenResponse, error) {
	return &v1.GetRoomInstanceAccessTokenResponse{}, nil
}
