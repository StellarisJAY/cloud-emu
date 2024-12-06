package biz

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/util"
	"github.com/bwmarrin/snowflake"
	"time"
)

type RoomInstanceUseCase struct {
	repo        RoomInstanceRepo
	snowflakeId *snowflake.Node
}

type RoomInstanceRepo interface {
	Create(ctx context.Context, roomInstance *RoomInstance) error
	Update(ctx context.Context, roomInstance *RoomInstance) error
	GetActiveInstanceByRoomId(ctx context.Context, roomId int64) (*RoomInstance, error)
	ListInstanceByRoomId(ctx context.Context, roomId int64, p *util.Pagination) ([]*RoomInstance, error)
}

type RoomInstance struct {
	RoomInstanceId int64     `json:"roomInstanceId"`
	RoomId         int64     `json:"roomId"`
	RoomName       string    `json:"roomName"`
	HostId         int64     `json:"hostId"`
	HostName       string    `json:"hostName"`
	EmulatorId     int64     `json:"emulatorId"`
	EmulatorName   string    `json:"emulatorName"`
	AddTime        time.Time `json:"addTime"`
	ServerUrl      string    `json:"serverUrl"`
	EndTime        time.Time `json:"endTime"`
	Status         int32     `json:"status"`
}

const (
	RoomInstanceStatusActive int32 = iota
	RoomInstanceStatusEnd
)

func NewRoomInstanceUseCase(repo RoomInstanceRepo, snowflakeId *snowflake.Node) *RoomInstanceUseCase {
	return &RoomInstanceUseCase{
		repo:        repo,
		snowflakeId: snowflakeId,
	}
}

func (uc *RoomInstanceUseCase) OpenRoomInstance(ctx context.Context, roomId int64) error {
	panic("implement me")
}

func (uc *RoomInstanceUseCase) ListRoomGameHistory(ctx context.Context, roomId int64, p *util.Pagination) ([]*RoomInstance, error) {
	instances, err := uc.repo.ListInstanceByRoomId(ctx, roomId, p)
	if err != nil {
		return nil, err
	}
	return instances, nil
}
