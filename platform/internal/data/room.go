package data

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
)

type RoomRepo struct {
	data *Data
}

const RoomTableName = "sys_room"
const RoomMemberTableName = "room_member"

func NewRoomRepo(data *Data) biz.RoomRepo {
	return &RoomRepo{data: data}
}

func (r *RoomRepo) Create(ctx context.Context, room *biz.Room) error {
	//TODO implement me
	panic("implement me")
}

func (r *RoomRepo) GetById(ctx context.Context, id int64) (*biz.Room, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoomRepo) Update(ctx context.Context, room *biz.Room) error {
	//TODO implement me
	panic("implement me")
}

func (r *RoomRepo) ListRooms(ctx context.Context, query biz.RoomQuery) ([]*biz.Room, error) {
	var rooms []*biz.Room
	err := r.data.db.Raw("SELECT sr.*, su.user_name AS 'host_name' "+
		"FROM sys_room sr "+
		"LEFT JOIN room_member rm ON sr.room_id = rm.room_id "+
		"INNER JOIN sys_user su ON su.user_id = sr.host_id "+
		"WHERE "+
		"rm.user_id = ? "+
		"OR "+
		"sr.host_id = ?", query.UserId, query.UserId).
		Find(&rooms).Debug().
		WithContext(ctx).
		Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
