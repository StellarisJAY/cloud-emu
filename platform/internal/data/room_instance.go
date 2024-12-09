package data

import (
	"context"
	"errors"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"gorm.io/gorm"
	"time"
)

type RoomInstanceRepo struct {
	data *Data
}

func NewRoomInstanceRepo(data *Data) biz.RoomInstanceRepo {
	return &RoomInstanceRepo{data: data}
}

const RoomInstanceTableName = "room_instance"

type RoomInstanceEntity struct {
	RoomInstanceId int64
	RoomId         int64
	AddTime        time.Time
	ServerUrl      string
	EmulatorId     int64
	EndTime        time.Time
	Status         int32
}

func (r *RoomInstanceRepo) Create(ctx context.Context, roomInstance *biz.RoomInstance) error {
	return r.data.db.Table(RoomInstanceTableName).
		Create(convertRoomInstanceBizToEntity(roomInstance)).
		WithContext(ctx).Error
}

func (r *RoomInstanceRepo) Update(ctx context.Context, roomInstance *biz.RoomInstance) error {
	return r.data.db.Table(RoomInstanceTableName).
		Where("room_instance_id=?", roomInstance.RoomInstanceId).
		UpdateColumns(convertRoomInstanceBizToEntity(roomInstance)).
		WithContext(ctx).
		Error
}

func (r *RoomInstanceRepo) GetActiveInstanceByRoomId(ctx context.Context, roomId int64) (*biz.RoomInstance, error) {
	var result *biz.RoomInstance
	err := r.data.db.Table(RoomInstanceTableName+" ri").
		Select("ri.*, e.emulator_name").
		Joins("LEFT JOIN emulator e ON ri.emulator_id = e.emulator_id").
		Where("ri.room_id = ?", roomId).
		Where("ri.status = ?", biz.RoomInstanceStatusActive).
		Scan(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *RoomInstanceRepo) ListInstanceByRoomId(ctx context.Context, roomId int64, p *common.Pagination) ([]*biz.RoomInstance, error) {
	var result []*biz.RoomInstance
	err := r.data.db.Table(RoomInstanceTableName+" ri").
		Select("ri.*", "emulator_name").
		Joins("LEFT JOIN emulator e ON ri.emulator_id = e.emulator_id").
		Where("ri.room_id = ?", roomId).
		Where("ri.status = ?", biz.RoomInstanceStatusActive).
		Scopes(common.WithPagination(p)).Debug().
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func convertRoomInstanceBizToEntity(dto *biz.RoomInstance) *RoomInstanceEntity {
	return &RoomInstanceEntity{
		RoomInstanceId: dto.RoomInstanceId,
		RoomId:         dto.RoomId,
		AddTime:        dto.AddTime,
		ServerUrl:      dto.ServerUrl,
		EmulatorId:     dto.EmulatorId,
		EndTime:        dto.EndTime,
		Status:         dto.Status,
	}
}
