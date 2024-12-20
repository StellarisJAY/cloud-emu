package service

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
)

type ButtonLayoutService struct {
	v1.UnimplementedButtonLayoutServer
	uc *biz.ButtonLayoutUseCase
}

func NewButtonLayoutService(uc *biz.ButtonLayoutUseCase) v1.ButtonLayoutServer {
	return &ButtonLayoutService{
		uc: uc,
	}
}

func (b *ButtonLayoutService) ListButtonLayout(ctx context.Context, request *v1.ListButtonLayoutRequest) (*v1.ListButtonLayoutResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b *ButtonLayoutService) CreateButtonLayout(ctx context.Context, request *v1.CreateButtonLayoutRequest) (*v1.CreateButtonLayoutResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b *ButtonLayoutService) UpdateButtonLayout(ctx context.Context, request *v1.UpdateButtonLayoutRequest) (*v1.UpdateButtonLayoutResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b *ButtonLayoutService) DeleteButtonLayout(ctx context.Context, request *v1.DeleteButtonLayoutRequest) (*v1.DeleteButtonLayoutResponse, error) {
	//TODO implement me
	panic("implement me")
}
