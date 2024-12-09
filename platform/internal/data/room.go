package data

import (
	"context"
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
	return r.data.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Table(RoomTableName).Create(convertRoomBizToEntity(room)).WithContext(ctx).Error
		if err != nil {
			return err
		}
		rm := &RoomMemberEntity{
			RoomMemberId: r.snowflakeId.Generate().Int64(),
			RoomId:       room.RoomId,
			UserId:       room.HostId,
			AddTime:      room.AddTime,
			Role:         biz.RoomMemberRoleHost,
		}
		return tx.Table(RoomMemberTableName).Create(rm).WithContext(ctx).Error
	})
}

func (r *RoomRepo) GetById(ctx context.Context, id int64) (*biz.Room, error) {
	room := &biz.Room{}
	err := r.data.db.Table(RoomTableName+"sr").
		Joins("LEFT JOIN room_member rm ON sr.room_id = rm.room_id ").
		Joins("INNER JOIN sys_user su ON su.user_id = sr.host_id ").
		Joins("LEFT JOIN room_instance ri ON sr.room_id = ri.room_id ").
		Where("sr.room_id = ?", id).
		WithContext(ctx).
		Scan(room).Error
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (r *RoomRepo) Update(ctx context.Context, room *biz.Room) error {
	//TODO implement me
	panic("implement me")
}

func (r *RoomRepo) ListRooms(ctx context.Context, query biz.RoomQuery, page *common.Pagination) ([]*biz.Room, error) {
	var rooms []*biz.Room
	db := r.data.db.Table(RoomTableName + " sr").
		Joins("LEFT JOIN room_member rm ON sr.room_id = rm.room_id ").
		Joins("INNER JOIN sys_user su ON su.user_id = sr.host_id ").
		Joins("LEFT JOIN room_instance ri ON sr.room_id = ri.room_id ")
	if query.HostId != 0 {
		db = db.Where("host_id = ?", query.HostId)
	}
	if query.MemberId != 0 {
		db = db.Where("rm.user_id = ?", query.MemberId)
	}
	if query.HostName != "" {
		db = db.Where("host_name LIKE %?%", query.HostName)
	}
	if query.RoomName != "" {
		db = db.Where("room_name LIKE %?%", query.RoomName)
	}
	if query.JoinType != 0 {
		db = db.Where("join_type = ?", query.JoinType)
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
