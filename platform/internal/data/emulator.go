package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/emulator"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/bwmarrin/snowflake"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type EmulatorRepo struct {
	data      *Data
	cache     *common.Cache[biz.Emulator]
	snowflake *snowflake.Node
}

const EmulatorTableName = "emulator"

func NewEmulatorRepo(data *Data, snowflake *snowflake.Node) biz.EmulatorRepo {
	e := &EmulatorRepo{
		data:      data,
		cache:     common.NewCache[biz.Emulator](data.redis),
		snowflake: snowflake,
	}
	if err := e.init(); err != nil {
		panic(err)
	}
	return e
}

func (e *EmulatorRepo) init() error {
	emulators := emulator.GetSupportedEmulators()
	if len(emulators) == 0 {
		panic("no supported emulator found")
	}
	for _, em := range emulators {
		emu := &biz.Emulator{
			EmulatorId:            e.snowflake.Generate().Int64(),
			EmulatorName:          em.Name,
			EmulatorType:          em.EmulatorType,
			Provider:              em.Provider,
			Description:           em.Description,
			SupportSave:           em.SupportSave,
			SupportGraphicSetting: em.SupportGraphicSettings,
		}
		err := e.data.db.Table(EmulatorTableName).Create(emu).Error
		if err != nil {
			var merr *mysql.MySQLError
			if errors.As(err, &merr) && merr.Number == 1062 {
				continue
			}
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				continue
			} else {
				return err
			}
		}
	}

	return nil
}

func (e *EmulatorRepo) cacheKey(id int64) string {
	return fmt.Sprintf("/emulator/%d", id)
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
	result, _ := e.cache.Get(ctx, e.cacheKey(id))
	if result != nil {
		return result, nil
	}
	err := e.data.DB(ctx).Table(EmulatorTableName).Where("emulator_id =?", id).WithContext(ctx).Scan(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	_ = e.cache.Set(ctx, e.cacheKey(id), result, 0)
	return result, nil
}

func (e *EmulatorRepo) GetByType(ctx context.Context, emulatorType string) (*biz.Emulator, error) {
	var result *biz.Emulator
	err := e.data.DB(ctx).Table(EmulatorTableName).Where("emulator_type =?", emulatorType).WithContext(ctx).Scan(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else {
		return result, err
	}
}
