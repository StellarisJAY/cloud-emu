package data

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
)

type ButtonLayoutRepo struct {
	data *Data
}

func NewButtonLayoutRepo(data *Data) biz.ButtonLayoutRepo {
	return &ButtonLayoutRepo{
		data: data,
	}
}

func (b *ButtonLayoutRepo) Create(ctx context.Context, buttonLayout *biz.ButtonLayout) error {
	//TODO implement me
	panic("implement me")
}

func (b *ButtonLayoutRepo) List(ctx context.Context, query biz.ButtonLayoutQuery, p *common.Pagination) ([]*biz.ButtonLayout, error) {
	//TODO implement me
	panic("implement me")
}

func (b *ButtonLayoutRepo) Update(ctx context.Context, buttonLayout *biz.ButtonLayout) error {
	//TODO implement me
	panic("implement me")
}

func (b *ButtonLayoutRepo) Delete(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}
