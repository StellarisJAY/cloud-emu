package data

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
)

type KeyboardBindingRepo struct {
	data *Data
}

func NewKeyboardBindingRepo(data *Data) biz.KeyboardBindingRepo {
	return &KeyboardBindingRepo{
		data: data,
	}
}

func (k *KeyboardBindingRepo) Create(ctx context.Context, keyboardBinding *biz.KeyboardBinding) error {
	//TODO implement me
	panic("implement me")
}

func (k *KeyboardBindingRepo) List(ctx context.Context, query biz.KeyboardBindingQuery, p *common.Pagination) ([]*biz.KeyboardBinding, error) {
	//TODO implement me
	panic("implement me")
}

func (k *KeyboardBindingRepo) Update(ctx context.Context, keyboardBinding *biz.KeyboardBinding) error {
	//TODO implement me
	panic("implement me")
}

func (k *KeyboardBindingRepo) Delete(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}
