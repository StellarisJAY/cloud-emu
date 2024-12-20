package service

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
)

type KeyboardBindingService struct {
	v1.UnimplementedKeyboardBindingServer
	uc *biz.KeyboardBindingUseCase
}

func NewKeyboardBindingService(uc *biz.KeyboardBindingUseCase) v1.KeyboardBindingServer {
	return &KeyboardBindingService{
		uc: uc,
	}
}

func (k *KeyboardBindingService) ListKeyboardBinding(ctx context.Context, request *v1.ListKeyboardBindingRequest) (*v1.ListKeyboardBindingResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k *KeyboardBindingService) CreateKeyboardBinding(ctx context.Context, request *v1.CreateKeyboardBindingRequest) (*v1.CreateKeyboardBindingResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k *KeyboardBindingService) UpdateKeyboardBinding(ctx context.Context, request *v1.UpdateKeyboardBindingRequest) (*v1.UpdateKeyboardBindingResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k *KeyboardBindingService) DeleteKeyboardBinding(ctx context.Context, request *v1.DeleteKeyboardBindingRequest) (*v1.DeleteKeyboardBindingResponse, error) {
	//TODO implement me
	panic("implement me")
}
