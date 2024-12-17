package data

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
)

type EmulatorRepo struct {
	data *Data
}

const EmulatorTableName = "emulator"

func NewEmulatorRepo(data *Data) biz.EmulatorRepo {
	return &EmulatorRepo{
		data: data,
	}
}

func (e *EmulatorRepo) ListEmulator(ctx context.Context, query biz.EmulatorQuery) ([]*biz.Emulator, error) {
	var result []*biz.Emulator
	d := e.data.DB(ctx).Table(EmulatorTableName)
	if query.EmulatorName != "" {
		d = d.Where("emulator_name LIKE ?", "%"+query.EmulatorName+"%")
	}
	if query.Provider != "" {
		d = d.Where("provider = ?", query.Provider)
	}
	err := d.WithContext(ctx).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
