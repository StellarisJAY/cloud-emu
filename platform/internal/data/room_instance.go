package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/hashicorp/consul/api"
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

func (r *RoomInstanceRepo) ListOnlineRoomMembers(ctx context.Context, roomInstance *biz.RoomInstance) ([]int64, error) {
	client, err := common.NewGRPCClient(roomInstance.ServerIp, int(roomInstance.RpcPort))
	if err != nil {
		return nil, err
	}
	defer client.Close()
	gameCli := v1.NewGameClient(client)
	resp, err := gameCli.ListOnlineRoomMember(ctx, &v1.ListOnlineRoomMemberRequest{
		RoomId: roomInstance.RoomId,
	})
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, errors.New(resp.Message)
	}
	return resp.RoomMemberIds, nil
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
