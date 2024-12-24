package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/hashicorp/consul/api"
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
	ServerIp       string
	EmulatorId     int64
	EndTime        time.Time
	Status         int32
	RpcPort        int32
	GameId         int64
}

func (r *RoomInstanceRepo) Create(ctx context.Context, roomInstance *biz.RoomInstance) error {
	return r.data.DB(ctx).Table(RoomInstanceTableName).
		Create(convertRoomInstanceBizToEntity(roomInstance)).
		WithContext(ctx).Error
}

func (r *RoomInstanceRepo) Update(ctx context.Context, roomInstance *biz.RoomInstance) error {
	return r.data.DB(ctx).Table(RoomInstanceTableName).
		Where("room_instance_id=?", roomInstance.RoomInstanceId).
		Updates(convertRoomInstanceBizToEntity(roomInstance)).
		WithContext(ctx).
		Error
}

func (r *RoomInstanceRepo) GetActiveInstanceByRoomId(ctx context.Context, roomId int64) (*biz.RoomInstance, error) {
	var result *biz.RoomInstance
	err := r.data.DB(ctx).Table(RoomInstanceTableName+" ri").
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
	err := r.data.DB(ctx).Table(RoomInstanceTableName+" ri").
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

func (r *RoomInstanceRepo) ListOnlineRoomMembers(ctx context.Context, roomInstance *biz.RoomInstance) ([]*biz.RoomMember, error) {
	// TODO 从游戏服务器获取房间实例在线成员
	return nil, nil
}

func (r *RoomInstanceRepo) SaveRoomInstance(_ context.Context, roomInstance *biz.RoomInstance) error {
	data, _ := json.Marshal(roomInstance)
	_, _, err := r.data.consul.KV().Acquire(&api.KVPair{
		Key:     fmt.Sprintf("cloudemu/room-instance/%d", roomInstance.RoomId),
		Value:   data,
		Session: roomInstance.SessionKey,
	}, nil)
	return err
}

func (r *RoomInstanceRepo) GetRoomInstance(_ context.Context, roomId int64) (*biz.RoomInstance, error) {
	res, _, err := r.data.consul.KV().Get(fmt.Sprintf("cloudemu/room-instance/%d", roomId), nil)
	if err != nil {
		return nil, err
	}
	if res == nil || len(res.Value) == 0 {
		return nil, nil
	}
	result := &biz.RoomInstance{}
	_ = json.Unmarshal(res.Value, result)
	return result, nil
}

func convertRoomInstanceBizToEntity(dto *biz.RoomInstance) *RoomInstanceEntity {
	return &RoomInstanceEntity{
		RoomInstanceId: dto.RoomInstanceId,
		RoomId:         dto.RoomId,
		AddTime:        dto.AddTime,
		ServerIp:       dto.ServerIp,
		EmulatorId:     dto.EmulatorId,
		EndTime:        dto.EndTime,
		Status:         dto.Status,
		RpcPort:        dto.RpcPort,
		GameId:         dto.GameId,
	}
}
