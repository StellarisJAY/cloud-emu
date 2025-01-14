package biz

import "context"

type GameServer struct {
	Address string
	Weight  int
}

// RoomMemberAuth 玩家连接服务器认证信息
type RoomMemberAuth struct {
	UserId int64
	Ip     string
	AppId  string
}

type RestartParams struct {
	UserId       int64
	EmulatorId   int64
	GameId       int64
	EmulatorCode string
	EmulatorType string
	GameName     string
	GameData     []byte
}

type LoadSaveParams struct {
	UserId       int64
	EmulatorId   int64
	GameId       int64
	EmulatorCode string
	EmulatorType string
	GameName     string
	GameData     []byte
	SaveData     []byte
}

type ControllerPlayer struct {
	ControllerId int32
	Label        string
	UserId       int64
}

type GameServerRepo interface {
	// ListActiveGameServers 服务发现 列出所有可用的游戏服务器
	ListActiveGameServers(ctx context.Context) ([]*GameServer, error)
	// OpenRoomInstance 在游戏服务器创建房间实例，返回创建实例用户的连接token和房间实例的心跳保活sessionKey
	OpenRoomInstance(ctx context.Context, instance *RoomInstance, auth RoomMemberAuth, gameData []byte) (string, string, error)
	// GetRoomInstanceToken 获取房间实例的连接token
	GetRoomInstanceToken(ctx context.Context, instance *RoomInstance, roomId int64, auth RoomMemberAuth) (string, error)
	// OpenGameConnection 创建连接，返回WebRTC的SDP Offer
	OpenGameConnection(ctx context.Context, instance *RoomInstance, token string, auth RoomMemberAuth) (string, error)
	// SdpAnswer 连接握手发送SDP Answer
	SdpAnswer(ctx context.Context, instance *RoomInstance, token string, auth RoomMemberAuth, sdpAnswer string) error
	// AddICECandidate 发送客户端ICE候选地址
	AddICECandidate(ctx context.Context, instance *RoomInstance, token string, auth RoomMemberAuth, candidate string) error
	// GetServerICECandidate 获取服务器ICE候选地址
	GetServerICECandidate(ctx context.Context, instance *RoomInstance, token string, auth RoomMemberAuth) ([]string, error)
	// RestartGameInstance 重启游戏实例
	RestartGameInstance(ctx context.Context, instance *RoomInstance, params RestartParams) error
	// SaveGame 保存游戏
	SaveGame(ctx context.Context, instance *RoomInstance, roomId, userId int64) (emulatorId, gameId int64, data []byte, err error)

	LoadSave(ctx context.Context, instance *RoomInstance, params LoadSaveParams) error

	GetControllerPlayers(ctx context.Context, instance *RoomInstance) ([]*ControllerPlayer, error)
}
