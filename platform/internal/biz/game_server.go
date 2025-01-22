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

type GraphicOptions struct {
	HighResolution bool
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
	// LoadSave 加载游戏存档
	LoadSave(ctx context.Context, instance *RoomInstance, params LoadSaveParams) error
	// GetControllerPlayers 获取控制器绑定的玩家
	GetControllerPlayers(ctx context.Context, instance *RoomInstance) ([]*ControllerPlayer, error)
	// SetControllerPlayer 设置控制器绑定的玩家
	SetControllerPlayer(context.Context, []*ControllerPlayer, *RoomInstance) error
	// GetGraphicOptions 获取游戏画面设置
	GetGraphicOptions(ctx context.Context, instance *RoomInstance) (*GraphicOptions, error)
	// SetGraphicOptions 设置游戏画面设置
	SetGraphicOptions(ctx context.Context, instance *RoomInstance, options *GraphicOptions) error
	// ApplyMacro 应用宏
	ApplyMacro(ctx context.Context, instance *RoomInstance, macro *Macro, userId int64) error
	// Shutdown 关闭游戏实例
	Shutdown(ctx context.Context, instance *RoomInstance) error

	GetEmulatorSpeed(ctx context.Context, instance *RoomInstance) (float64, error)
	SetEmulatorSpeed(ctx context.Context, instance *RoomInstance, boost float64) (float64, error)
}
