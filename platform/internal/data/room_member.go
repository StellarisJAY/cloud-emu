package data

import "time"

type RoomMemberEntity struct {
	RoomMemberId int64
	RoomId       int64
	UserId       int64
	AddTime      time.Time
	Role         int32
}

const RoomMemberTableName = "room_member"
