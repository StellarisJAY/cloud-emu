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
	return r.data.DB(ctx).Table(RoomMemberTableName).WithContext(ctx).Create(convertRoomMemberBizToEntity(member)).Error
}

func (r *RoomMemberRepo) Update(ctx context.Context, member *biz.RoomMember) error {
	return r.data.DB(ctx).Table(RoomMemberTableName).Where("room_member_id=?", member.RoomMemberId).WithContext(ctx).Updates(convertRoomMemberBizToEntity(member)).Error
}

func (r *RoomMemberRepo) List(ctx context.Context, roomId int64) ([]*biz.RoomMember, error) {
	var members []*biz.RoomMember
	err := r.data.DB(ctx).Table(RoomMemberTableName+" rm").
		Joins("INNER JOIN "+UserTableName+" su ON su.user_id = rm.user_id").
		Select("rm.*, su.user_name, su.nick_name").
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
	err := r.data.DB(ctx).Table(RoomMemberTableName+" rm").
		Joins("INNER JOIN "+UserTableName+" su ON su.user_id = rm.user_id").
		Select("rm.*, su.user_name, su.nick_name").
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

func convertRoomMemberBizToEntity(member *biz.RoomMember) *RoomMemberEntity {
	return &RoomMemberEntity{
		RoomMemberId: member.RoomMemberId,
		RoomId:       member.RoomId,
		UserId:       member.UserId,
		AddTime:      time.Now(),
		Role:         member.Role,
	}
}

func (r *RoomMemberRepo) CountRoomMember(ctx context.Context, roomId int64) (int32, error) {
	var count int64 = 0
	err := r.data.DB(ctx).Table(RoomMemberTableName).
		Where("room_id = ?", roomId).
		WithContext(ctx).
		Count(&count).
		Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}
