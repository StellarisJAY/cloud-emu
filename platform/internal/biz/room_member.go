package biz

import (
	"context"
	"errors"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/bwmarrin/snowflake"
	"slices"
	"time"
)

const (
	RoomMemberRoleHost int32 = iota + 1
	RoomMemberRolePlayer
	RoomMemberRoleObserver
)

const (
	RoomMemberStatusJoined int32 = iota + 1
	RoomMemberStatusInvited
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
}

type RoomMemberUseCase struct {
	roomMemberRepo   RoomMemberRepo
	roomRepo         RoomRepo
	roomInstanceRepo RoomInstanceRepo
	notificationRepo NotificationRepo
	snowflakeId      *snowflake.Node
	userRepo         UserRepo
	tm               Transaction
}

func NewRoomMemberUseCase(roomMemberRepo RoomMemberRepo, roomRepo RoomRepo, roomInstanceRepo RoomInstanceRepo,
	notificationRepo NotificationRepo, snowflakeId *snowflake.Node, userRepo UserRepo, tm Transaction) *RoomMemberUseCase {
	return &RoomMemberUseCase{
		roomMemberRepo:   roomMemberRepo,
		roomRepo:         roomRepo,
		roomInstanceRepo: roomInstanceRepo,
		notificationRepo: notificationRepo,
		snowflakeId:      snowflakeId,
		userRepo:         userRepo,
		tm:               tm,
	}
}

func (uc *RoomMemberUseCase) ListRoomMembers(ctx context.Context, roomId int64) ([]*RoomMember, error) {
	members, err := uc.roomMemberRepo.List(ctx, roomId)
	if err != nil {
		return nil, err
	}
	roomInstance, _ := uc.roomInstanceRepo.GetActiveInstanceByRoomId(ctx, roomId)
	// 没有活跃房间实例，成员均不在线
	if roomInstance == nil {
		return members, nil
	}
	// 获取在线成员列表
	onlineMembers, _ := uc.roomInstanceRepo.ListOnlineRoomMembers(ctx, roomInstance)
	if len(onlineMembers) > 0 {
		for _, member := range onlineMembers {
			idx := slices.IndexFunc(members, func(item *RoomMember) bool {
				return item.RoomMemberId == member.RoomMemberId
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
		member, _ := uc.roomMemberRepo.GetByRoomAndUser(ctx, roomId, invitedUserId)
		if member != nil {
			return errors.New("无法邀请该用户")
		}
		user, _ := uc.roomMemberRepo.GetByRoomAndUser(ctx, roomId, userId)
		if user == nil || user.Role != RoomMemberRoleHost {
			return v1.ErrorAccessDenied("无法邀请该用户")
		}

		if u, _ := uc.userRepo.GetById(ctx, invitedUserId); u == nil || u.Status != UserStatusAvailable {
			return errors.New("无法邀请该用户")
		}

		member = &RoomMember{
			RoomMemberId: uc.snowflakeId.Generate().Int64(),
			RoomId:       roomId,
			UserId:       invitedUserId,
			Role:         RoomMemberRolePlayer,
			AddTime:      time.Now().Local(),
			Status:       RoomMemberStatusInvited,
		}
		err := uc.roomMemberRepo.Create(ctx, member)
		if err != nil {
			return errors.New("无法邀请该用户")
		}
		notice := Notification{
			NotificationId: uc.snowflakeId.Generate().Int64(),
			Type:           NotificationTypeInvitation,
			SenderId:       userId,
			Content:        "",
			AddTime:        time.Now().Local(),
			ReceiverId:     invitedUserId,
		}
		if err = uc.notificationRepo.Create(ctx, &notice); err != nil {
			return errors.New("无法邀请该用户")
		}
		return nil
	})
}
