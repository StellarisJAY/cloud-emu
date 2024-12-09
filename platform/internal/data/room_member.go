package data

import (
	"context"
	"errors"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"gorm.io/gorm"
	"time"
)

type RoomMemberEntity struct {
	RoomMemberId int64
	RoomId       int64
	UserId       int64
	AddTime      time.Time
	Role         int32
}

const RoomMemberTableName = "room_member"

type RoomMemberRepo struct {
	data *Data
}

func NewRoomMemberRepo(data *Data) biz.RoomMemberRepo {
	return &RoomMemberRepo{data: data}
}

func (r *RoomMemberRepo) Create(ctx context.Context, member *biz.RoomMember) error {
	return r.data.db.Table(RoomMemberTableName).WithContext(ctx).Create(member).Error
}

func (r *RoomMemberRepo) Update(ctx context.Context, member *biz.RoomMember) error {
	return r.data.db.Table(RoomMemberTableName).Where("room_member_id=?", member.RoomMemberId).WithContext(ctx).Updates(member).Error
}

func (r *RoomMemberRepo) List(ctx context.Context, roomId int64) ([]*biz.RoomMember, error) {
	var members []*biz.RoomMember
	err := r.data.db.Table(RoomMemberTableName+" rm").
		Joins("INNER JOIN "+UserTableName+" su ON su.user_id = rm.user_id").
		Where("rm.room_id=?", roomId).
		WithContext(ctx).
		Scan(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *RoomMemberRepo) GetByRoomAndUser(ctx context.Context, roomId, userId int64) (*biz.RoomMember, error) {
	var member *biz.RoomMember
	err := r.data.db.Table(RoomMemberTableName+" rm").
		Joins("INNER JOIN "+UserTableName+" su ON su.user_id = rm.user_id").
		Where("rm.room_id=?", roomId).
		Where("rm.user_id=?", userId).
		WithContext(ctx).
		Scan(&member).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return member, nil
}
