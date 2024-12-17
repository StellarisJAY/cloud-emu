package biz

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/bwmarrin/snowflake"
	"github.com/go-redsync/redsync/v4"
	"time"
)

type Room struct {
	RoomId       int64     `json:"roomId"`
	RoomName     string    `json:"roomName"`
	Description  string    `json:"description"`
	HostId       int64     `json:"hostId"`
	HostName     string    `json:"hostName"`
	MemberCount  int32     `json:"memberCount"`
	MemberLimit  int32     `json:"memberLimit"`
	AddTime      time.Time `json:"addTime"`
	Password     string    `json:"password"`
	JoinType     int32     `json:"joinRule"`
	EmulatorId   int64     `json:"emulatorId"`
	EmulatorName string    `json:"emulatorName"`
}

type RoomUseCase struct {
	repo           RoomRepo
	snowflakeId    *snowflake.Node
	userRepo       UserRepo
	tm             Transaction
	roomMemberRepo RoomMemberRepo
	redSync        *redsync.Redsync
}

type RoomQuery struct {
	HostId     int64
	RoomName   string
	HostName   string
	JoinType   int32
	MemberId   int64
	EmulatorId int64
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

func NewRoomUseCase(repo RoomRepo, snowflakeId *snowflake.Node, userRepo UserRepo, tm Transaction,
	roomMemberRepo RoomMemberRepo, redsync *redsync.Redsync) *RoomUseCase {
	return &RoomUseCase{repo: repo, snowflakeId: snowflakeId, userRepo: userRepo, tm: tm, roomMemberRepo: roomMemberRepo,
		redSync: redsync}
}

func (r *RoomUseCase) Create(ctx context.Context, room *Room) error {
	return r.tm.Tx(ctx, func(ctx context.Context) error {
		user, _ := r.userRepo.GetById(ctx, room.HostId)
		if user == nil || user.Status != UserStatusAvailable {
			return v1.ErrorAccessDenied("当前用户无法创建房间")
		}
		room.RoomId = r.snowflakeId.Generate().Int64()
		if room.Password != "" {
			hash := sha256.Sum256([]byte(room.Password))
			room.Password = hex.EncodeToString(hash[:])
		}
		room.AddTime = time.Now().Local()
		if err := r.repo.Create(ctx, room); err != nil {
			return err
		}
		return r.roomMemberRepo.Create(ctx, &RoomMember{
			RoomMemberId: r.snowflakeId.Generate().Int64(),
			RoomId:       room.RoomId,
			UserId:       room.HostId,
			Role:         RoomMemberRoleHost,
			AddTime:      room.AddTime,
			Status:       RoomMemberStatusJoined,
		})
	})
}

func (r *RoomUseCase) ListMyRooms(ctx context.Context, userId int64, query RoomQuery, page *common.Pagination) ([]*Room, error) {
	query.MemberId = userId
	rooms, err := r.repo.ListRooms(ctx, query, page)
	if err != nil {
		return nil, err
	}
	for _, room := range rooms {
		if err := r.buildRoomDto(ctx, room); err != nil {
			return nil, err
		}
	}
	return rooms, nil
}

func (r *RoomUseCase) buildRoomDto(ctx context.Context, room *Room) error {
	count, err := r.roomMemberRepo.CountRoomMember(ctx, room.RoomId)
	if err != nil {
		return err
	}
	room.MemberCount = count
	// TODO 获取模拟器信息
	return nil
}
