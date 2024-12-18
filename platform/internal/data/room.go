package data

import (
	"context"
	"errors"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
	"time"
)

type RoomRepo struct {
	data        *Data
	snowflakeId *snowflake.Node
}

const RoomTableName = "sys_room"

func NewRoomRepo(data *Data, snowflakeId *snowflake.Node) biz.RoomRepo {
	return &RoomRepo{data: data, snowflakeId: snowflakeId}
}

type RoomEntity struct {
	RoomId      int64     `json:"roomId"`
	RoomName    string    `json:"roomName"`
	Description string    `json:"description"`
	HostId      int64     `json:"hostId"`
	MemberLimit int32     `json:"memberLimit"`
	AddTime     time.Time `json:"addTime"`
	Password    string    `json:"password"`
	JoinType    int32     `json:"joinRule"`
}

func (r *RoomRepo) Create(ctx context.Context, room *biz.Room) error {
	return r.data.DB(ctx).Table(RoomTableName).Create(convertRoomBizToEntity(room)).Error
}

func (r *RoomRepo) GetById(ctx context.Context, id int64) (*biz.Room, error) {
	var room *biz.Room
	err := r.data.DB(ctx).Table(RoomTableName+" sr").Select("sr.*, su.user_name AS host_name, ri.game_id, eg.game_name").
		Joins("LEFT JOIN room_member rm ON sr.room_id = rm.room_id ").
		Joins("INNER JOIN sys_user su ON su.user_id = sr.host_id ").
		Joins("LEFT JOIN room_instance ri ON sr.room_id = ri.room_id ").
		Joins("LEFT JOIN emulator_game eg ON eg.game_id = ri.game_id").
		Where("sr.room_id = ?", id).
		WithContext(ctx).
		Scan(&room).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (r *RoomRepo) Update(ctx context.Context, room *biz.Room) error {
	return r.data.DB(ctx).Table(RoomTableName).
		Where("room_id = ?", room.RoomId).
		Updates(convertRoomBizToEntity(room)).
		WithContext(ctx).
		Error
}

func (r *RoomRepo) ListRooms(ctx context.Context, query biz.RoomQuery, page *common.Pagination) ([]*biz.Room, error) {
	var rooms []*biz.Room
	db := r.data.DB(ctx).Table(RoomTableName + " sr").Select("sr.*, su.user_name AS host_name").
		Joins("INNER JOIN room_member rm ON sr.room_id = rm.room_id ").
		Joins("INNER JOIN sys_user su ON su.user_id = sr.host_id ").
		Joins("LEFT JOIN room_instance ri ON sr.room_id = ri.room_id ")
	if query.MemberId != 0 {
		db = db.Where("rm.user_id = ?", query.MemberId)
	}
	if query.HostName != "" {
		db = db.Where("su.user_name LIKE ?", "%"+query.HostName+"%")
	}
	if query.RoomName != "" {
		db = db.Where("room_name LIKE ?", "%"+query.RoomName+"%")
	}
	if query.JoinType != 0 {
		db = db.Where("join_type = ?", query.JoinType)
	}
	if query.EmulatorId != 0 {
		db = db.Where("emulator_id = ?", query.EmulatorId)
	}
	err := db.Scopes(common.WithPagination(page)).WithContext(ctx).Scan(&rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func convertRoomBizToEntity(room *biz.Room) *RoomEntity {
	return &RoomEntity{
		RoomId:      room.RoomId,
		RoomName:    room.RoomName,
		Description: room.Description,
		HostId:      room.HostId,
		MemberLimit: room.MemberLimit,
		AddTime:     room.AddTime,
		Password:    room.Password,
		JoinType:    room.JoinType,
	}
}
