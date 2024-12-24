package biz

import "context"

type GameServer struct {
	Address string
	Weight  int
}

type RoomMemberAuth struct {
	UserId int64
	Ip     string
	AppId  string
}

type GameServerRepo interface {
	ListActiveGameServers(ctx context.Context) ([]*GameServer, error)
	OpenRoomInstance(ctx context.Context, instance *RoomInstance, auth RoomMemberAuth) (string, string, error)
	GetRoomInstanceToken(ctx context.Context, instance *RoomInstance, roomId int64, auth RoomMemberAuth) (string, error)
	OpenGameConnection(ctx context.Context, instance *RoomInstance, token string, auth RoomMemberAuth) (string, error)
	SdpAnswer(ctx context.Context, instance *RoomInstance, token string, auth RoomMemberAuth, sdpAnswer string) error
	AddICECandidate(ctx context.Context, instance *RoomInstance, token string, auth RoomMemberAuth, candidate string) error
	GetServerICECandidate(ctx context.Context, instance *RoomInstance, token string, auth RoomMemberAuth) ([]string, error)
	RestartGameInstance(ctx context.Context, instance *RoomInstance, userId int64, emulatorType, gameName, gameUrl string) error
}
