package biz

import (
	"context"
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
}

type RoomMemberRepo interface {
	Create(ctx context.Context, member *RoomMember) error
	Update(ctx context.Context, member *RoomMember) error
	List(ctx context.Context, roomId int64) ([]*RoomMember, error)
	GetByRoomAndUser(ctx context.Context, roomId, userId int64) (*RoomMember, error)
}
