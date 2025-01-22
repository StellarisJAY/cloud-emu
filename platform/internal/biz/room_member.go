package biz

import (
	"context"
	"encoding/json"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/bwmarrin/snowflake"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"slices"
	"strconv"
	"time"
)

const (
	RoomMemberRoleHost int32 = iota + 1
	RoomMemberRolePlayer
	RoomMemberRoleObserver
)

const (
	RoomMemberStatusJoined int32 = iota + 1
	RoomMemberStatusBanned
)

type RoomMember struct {
	RoomMemberId int64
	RoomId       int64
	UserId       int64
	Role         int32
	AddTime      time.Time
	UserName     string
	NickName     string
	Online       bool
	Status       int32
}

type RoomMemberRepo interface {
	Create(ctx context.Context, member *RoomMember) error
	Update(ctx context.Context, member *RoomMember) error
	List(ctx context.Context, roomId int64) ([]*RoomMember, error)
	GetByRoomAndUser(ctx context.Context, roomId, userId int64) (*RoomMember, error)
	CountRoomMember(ctx context.Context, roomId int64) (int32, error)
	DeleteRoomAllMembers(ctx context.Context, roomId int64) error
}

type RoomMemberUseCase struct {
	roomMemberRepo   RoomMemberRepo
	roomRepo         RoomRepo
	roomInstanceRepo RoomInstanceRepo
	notificationRepo NotificationRepo
	snowflakeId      *snowflake.Node
	userRepo         UserRepo
	tm               Transaction
	logger           *log.Helper
}

func NewRoomMemberUseCase(roomMemberRepo RoomMemberRepo, roomRepo RoomRepo, roomInstanceRepo RoomInstanceRepo,
	notificationRepo NotificationRepo, snowflakeId *snowflake.Node, userRepo UserRepo, tm Transaction, logger log.Logger) *RoomMemberUseCase {
	return &RoomMemberUseCase{
		roomMemberRepo:   roomMemberRepo,
		roomRepo:         roomRepo,
		roomInstanceRepo: roomInstanceRepo,
		notificationRepo: notificationRepo,
		snowflakeId:      snowflakeId,
		userRepo:         userRepo,
		tm:               tm,
		logger:           log.NewHelper(logger),
	}
}

func (uc *RoomMemberUseCase) ListRoomMembers(ctx context.Context, roomId int64) ([]*RoomMember, error) {
	members, err := uc.roomMemberRepo.List(ctx, roomId)
	if err != nil {
		return nil, err
	}
	roomInstance, _ := uc.roomInstanceRepo.GetRoomInstance(ctx, roomId)
	// 没有活跃房间实例，成员均不在线
	if roomInstance == nil {
		return members, nil
	}
	// 获取在线成员列表
	onlineMembers, _ := uc.roomInstanceRepo.ListOnlineRoomMembers(ctx, roomInstance)
	if len(onlineMembers) > 0 {
		for _, uid := range onlineMembers {
			idx := slices.IndexFunc(members, func(item *RoomMember) bool {
				return item.UserId == uid
			})
			if idx >= 0 {
				members[idx].Online = true
			}
		}
	}
	return members, nil
}

func (uc *RoomMemberUseCase) InviteRoomMember(ctx context.Context, userId int64, invitedUserId int64, roomId int64) error {
	return uc.tm.Tx(ctx, func(ctx context.Context) error {
		room, _ := uc.roomRepo.GetById(ctx, roomId)
		if room == nil || room.HostId != userId {
			return v1.ErrorAccessDenied("无法邀请用户加入房间")
		}
		member, _ := uc.roomMemberRepo.GetByRoomAndUser(ctx, roomId, invitedUserId)
		if member != nil {
			return errors.New(500, "", "无法邀请该用户")
		}
		if u, _ := uc.userRepo.GetById(ctx, invitedUserId); u == nil || u.Status != UserStatusAvailable {
			return errors.New(500, "", "无法邀请该用户")
		}
		content, _ := json.Marshal(struct {
			RoomId   string `json:"roomId"`
			RoomName string `json:"roomName"`
		}{strconv.FormatInt(roomId, 10), room.RoomName})
		notice := Notification{
			NotificationId: uc.snowflakeId.Generate().Int64(),
			Type:           NotificationTypeInvitation,
			SenderId:       userId,
			Content:        string(content),
			AddTime:        time.Now().Local(),
			ReceiverId:     invitedUserId,
		}
		if err := uc.notificationRepo.Create(ctx, &notice); err != nil {
			return errors.New(500, "", "无法邀请该用户")
		}
		return nil
	})
}

func (uc *RoomMemberUseCase) GetByRoomAndUser(ctx context.Context, roomId, userId int64) (*RoomMember, error) {
	rm, err := uc.roomMemberRepo.GetByRoomAndUser(ctx, roomId, userId)
	if err != nil {
		return nil, errors.New(500, "Database Error", "查询出错")
	}
	if rm == nil {
		return nil, v1.ErrorNotFound("用户不存在")
	}
	return rm, nil
}
