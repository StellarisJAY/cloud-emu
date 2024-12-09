package biz

import "context"

type GameServer struct {
	Address string
	Weight  int
}

type GameServerRepo interface {
	ListActiveGameServers(ctx context.Context) ([]*GameServer, error)
	OpenRoomInstance(ctx context.Context, instance *RoomInstance) (string, error)
	GetRoomInstanceToken(ctx context.Context, instance *RoomInstance, roomId, userId int64) (string, error)
}
