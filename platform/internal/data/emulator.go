package data

import (
	"context"
	"errors"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"gorm.io/gorm"
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

func (e *EmulatorRepo) GetById(ctx context.Context, id int64) (*biz.Emulator, error) {
	var result *biz.Emulator
	err := e.data.DB(ctx).Table(EmulatorTableName).Where("emulator_id =?", id).WithContext(ctx).Scan(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else {
		return result, err
	}
}
