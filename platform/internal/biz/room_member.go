package biz

import (
	"context"
	"slices"
	"time"
)

const (
	RoomMemberRoleHost int32 = iota + 1
	RoomMemberRolePlayer
	RoomMemberRoleObserver
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
}

type RoomMemberRepo interface {
	Create(ctx context.Context, member *RoomMember) error
	Update(ctx context.Context, member *RoomMember) error
	List(ctx context.Context, roomId int64) ([]*RoomMember, error)
	GetByRoomAndUser(ctx context.Context, roomId, userId int64) (*RoomMember, error)
}

type RoomMemberUseCase struct {
	roomMemberRepo   RoomMemberRepo
	roomRepo         RoomRepo
	roomInstanceRepo RoomInstanceRepo
}

func NewRoomMemberUseCase(roomMemberRepo RoomMemberRepo, roomRepo RoomRepo, roomInstanceRepo RoomInstanceRepo) *RoomMemberUseCase {
	return &RoomMemberUseCase{
		roomMemberRepo:   roomMemberRepo,
		roomRepo:         roomRepo,
		roomInstanceRepo: roomInstanceRepo,
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
