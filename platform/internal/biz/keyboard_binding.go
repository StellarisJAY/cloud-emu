package biz

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/go-kratos/kratos/v2/log"
)

type KeyboardBinding struct {
	BindingId   int64
	BindingName string
	AddUser     int64
	EmulatorId  int64
	Binding     string
}

type KeyboardBindingQuery struct {
	BindingName string
	AddUser     int64
	EmulatorId  int64
}

type KeyboardBindingRepo interface {
	Create(ctx context.Context, keyboardBinding *KeyboardBinding) error
	List(ctx context.Context, query KeyboardBindingQuery, p *common.Pagination) ([]*KeyboardBinding, error)
	Update(ctx context.Context, keyboardBinding *KeyboardBinding) error
	Delete(ctx context.Context, id int64) error
}

type KeyboardBindingUseCase struct {
	repo   KeyboardBindingRepo
	logger *log.Helper
}

func NewKeyboardBindingUseCase(repo KeyboardBindingRepo, logger log.Logger) *KeyboardBindingUseCase {
	return &KeyboardBindingUseCase{repo: repo, logger: log.NewHelper(logger)}
}

func (uc *KeyboardBindingUseCase) Create(ctx context.Context, keyboardBinding *KeyboardBinding) error {
	panic("implement me")
}

func (uc *KeyboardBindingUseCase) List(ctx context.Context, query KeyboardBindingQuery, p *common.Pagination) ([]*KeyboardBinding, error) {
	panic("implement me")
}

func (uc *KeyboardBindingUseCase) Update(ctx context.Context, keyboardBinding *KeyboardBinding) error {
	panic("implement me")
}

func (uc *KeyboardBindingUseCase) Delete(ctx context.Context, id int64) error {
	panic("implement me")
}
