package biz

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"time"
)

type Room struct {
	RoomId      int64     `json:"roomId"`
	RoomName    string    `json:"roomName"`
	Description string    `json:"description"`
	HostId      int64     `json:"hostId"`
	HostName    string    `json:"hostName"`
	MemberLimit int32     `json:"memberLimit"`
	AddTime     time.Time `json:"addTime"`
	Password    string    `json:"password"`
	JoinType    int32     `json:"joinRule"`
}

type RoomUseCase struct {
	repo        RoomRepo
	snowflakeId *snowflake.Node
}

type RoomQuery struct {
	UserId   int64
	RoomName string
	HostName string
	JoinRule int32
	page     int32
	pageSize int32
}

type RoomRepo interface {
	Create(ctx context.Context, room *Room) error
	GetById(ctx context.Context, id int64) (*Room, error)
	Update(ctx context.Context, room *Room) error
	ListRooms(ctx context.Context, query RoomQuery) ([]*Room, error)
}

func NewRoomUseCase(repo RoomRepo, snowflakeId *snowflake.Node) *RoomUseCase {
	return &RoomUseCase{repo: repo, snowflakeId: snowflakeId}
}

func (uc *RoomUseCase) ListMyRooms(ctx context.Context, userId int64) ([]*Room, error) {
	query := RoomQuery{UserId: userId}
	rooms, err := uc.repo.ListRooms(ctx, query)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
