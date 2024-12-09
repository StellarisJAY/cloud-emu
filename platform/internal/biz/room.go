package biz

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/bwmarrin/snowflake"
	"time"
)

type Room struct {
	RoomId       int64     `json:"roomId"`
	RoomName     string    `json:"roomName"`
	Description  string    `json:"description"`
	HostId       int64     `json:"hostId"`
	HostName     string    `json:"hostName"`
	MemberLimit  int32     `json:"memberLimit"`
	AddTime      time.Time `json:"addTime"`
	Password     string    `json:"password"`
	JoinType     int32     `json:"joinRule"`
	EmulatorId   int64     `json:"emulatorId"`
	EmulatorName string    `json:"emulatorName"`
}

type RoomUseCase struct {
	repo        RoomRepo
	snowflakeId *snowflake.Node
}

type RoomQuery struct {
	HostId   int64
	RoomName string
	HostName string
	JoinType int32
	MemberId int64
}

const (
	RoomJoinTypePublic int32 = iota + 1
	RoomJoinTypePassword
	RoomJoinTypeInvite
)

type RoomRepo interface {
	Create(ctx context.Context, room *Room) error
	GetById(ctx context.Context, id int64) (*Room, error)
	Update(ctx context.Context, room *Room) error
	ListRooms(ctx context.Context, query RoomQuery, page *common.Pagination) ([]*Room, error)
}

func NewRoomUseCase(repo RoomRepo, snowflakeId *snowflake.Node) *RoomUseCase {
	return &RoomUseCase{repo: repo, snowflakeId: snowflakeId}
}

func (uc *RoomUseCase) Create(ctx context.Context, room *Room) error {
	room.RoomId = uc.snowflakeId.Generate().Int64()
	if room.Password != "" {
		hash := sha256.Sum256([]byte(room.Password))
		room.Password = hex.EncodeToString(hash[:])
	}
	room.AddTime = time.Now().Local()
	return uc.repo.Create(ctx, room)
}

func (uc *RoomUseCase) ListMyRooms(ctx context.Context, userId int64, query RoomQuery, page *common.Pagination) ([]*Room, error) {
	query.MemberId = userId
	query.HostId = 0
	rooms, err := uc.repo.ListRooms(ctx, query, page)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
