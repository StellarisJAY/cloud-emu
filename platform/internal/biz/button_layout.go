package biz

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/go-kratos/kratos/v2/log"
)

type ButtonLayout struct {
	LayoutId   int64
	LayoutName string
	AddUser    int64
	EmulatorId int64
	Layout     string
}

type ButtonLayoutQuery struct {
	LayoutName string
	AddUser    int64
	EmulatorId int64
}

type ButtonLayoutRepo interface {
	Create(ctx context.Context, buttonLayout *ButtonLayout) error
	List(ctx context.Context, query ButtonLayoutQuery, p *common.Pagination) ([]*ButtonLayout, error)
	Update(ctx context.Context, buttonLayout *ButtonLayout) error
	Delete(ctx context.Context, id int64) error
}

type ButtonLayoutUseCase struct {
	repo   ButtonLayoutRepo
	logger *log.Helper
}

func NewButtonLayoutUseCase(repo ButtonLayoutRepo, logger log.Logger) *ButtonLayoutUseCase {
	return &ButtonLayoutUseCase{repo: repo, logger: log.NewHelper(logger)}
}

func (uc *ButtonLayoutUseCase) Create(ctx context.Context, buttonLayout *ButtonLayout) error {
	panic("implement me")
}

func (uc *ButtonLayoutUseCase) List(ctx context.Context, query ButtonLayoutQuery, p *common.Pagination) ([]*ButtonLayout, error) {
	panic("implement me")
}

func (uc *ButtonLayoutUseCase) Update(ctx context.Context, buttonLayout *ButtonLayout) error {
	panic("implement me")
}

func (uc *ButtonLayoutUseCase) Delete(ctx context.Context, id int64) error {
	panic("implement me")
}
