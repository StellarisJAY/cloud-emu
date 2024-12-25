package biz

import (
	"context"
	"encoding/json"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/pion/webrtc/v3"
	"sync"
	"time"
)

const (
	DefaultAudioSampleRate = 48000
)

type GameSave struct {
	Id         int64  `json:"id"`
	RoomId     int64  `json:"roomId"`
	Game       string `json:"game"`
	Data       []byte `json:"data"`
	CreateTime int64  `json:"createTime"`
	ExitSave   bool   `json:"exitSave"`
}
type GameInstanceStats struct {
	RoomId            int64         `json:"roomId"`
	Connections       int           `json:"connections"`
	ActiveConnections int           `json:"activeConnections"`
	Game              string        `json:"game"`
	Uptime            time.Duration `json:"uptime"`
}

type GraphicOptions struct {
	HighResOpen  bool `json:"highResOpen"`
	ReverseColor bool `json:"reverseColor"`
	Grayscale    bool `json:"grayscale"`
}

type GameFileRepo interface {
	GetGameData(ctx context.Context, game string) ([]byte, error)
	GetSavedGame(ctx context.Context, id int64) (*GameSave, error)
	SaveGame(ctx context.Context, save *GameSave) error
	ListSaves(ctx context.Context, roomId int64, page, pageSize int32) ([]*GameSave, int32, error)
	DeleteSave(ctx context.Context, saveId int64) error
	GetExitSave(ctx context.Context, roomId int64) (*GameSave, error)
}

type GameServerUseCase struct {
	gameFileRepo   GameFileRepo
	memberAuthRepo MemberAuthRepo
	gameInstances  map[int64]*GameInstance
	mutex          *sync.RWMutex
	logger         *log.Helper
	consul         *api.Client
}

type CreateRoomInstanceParams struct {
	RoomId         int64
	Auth           *MemberAuthInfo
	RoomInstanceId int64
}

type RestartParams struct {
	RoomId       int64
	UserId       int64
	EmulatorType string
	GameName     string
	GameUrl      string
	EmulatorId   int64
	GameId       int64
}

func NewGameServerUseCase(gameFileRepo GameFileRepo, memberAuthRepo MemberAuthRepo, logger log.Logger, consul *api.Client) *GameServerUseCase {
	return &GameServerUseCase{
		gameFileRepo:   gameFileRepo,
		memberAuthRepo: memberAuthRepo,
		logger:         log.NewHelper(logger),
		mutex:          &sync.RWMutex{},
		gameInstances:  make(map[int64]*GameInstance),
		consul:         consul,
	}
}

func (uc *GameServerUseCase) CreateRoomInstance(ctx context.Context, params CreateRoomInstanceParams) (string, string, error) {
	instance, err := makeGameInstance(params.RoomId)
	if err != nil {
		return "", "", v1.ErrorServiceError("创建游戏实例出错")
	}
	uid, _ := uuid.NewUUID()
	token := uid.String()
	if err := uc.memberAuthRepo.StoreAuthInfo(token, params.RoomId, params.Auth); err != nil {
		return "", "", v1.ErrorServiceError("创建玩家token出错")
	}
	consumerCtx, consumerCancel := context.WithCancel(context.Background())
	go instance.messageConsumer(consumerCtx)
	instance.cancel = consumerCancel
	uc.mutex.Lock()
	defer uc.mutex.Unlock()
	uc.gameInstances[params.RoomId] = instance
	// 创建consul session
	sessionKey, err := uc.createGameInstanceSession(instance)
	return token, sessionKey, err
}

func (uc *GameServerUseCase) GetRoomInstanceToken(ctx context.Context, roomId int64, auth *MemberAuthInfo) (string, error) {
	uc.mutex.RLock()
	_, ok := uc.gameInstances[roomId]
	uc.mutex.RUnlock()
	if !ok {
		return "", v1.ErrorAccessDenied("连接失败，游戏实例不存在")
	}
	uid, _ := uuid.NewUUID()
	token := uid.String()
	if err := uc.memberAuthRepo.StoreAuthInfo(token, roomId, auth); err != nil {
		return "", v1.ErrorServiceError("创建玩家token出错")
	}
	return token, nil
}

// OpenGameConnection 创建游戏连接，检查用户token，创建游戏实例webrtc连接，返回sdp offer
func (uc *GameServerUseCase) OpenGameConnection(_ context.Context, roomId int64, token string, auth *MemberAuthInfo) (string, error) {
	instance, err := uc.getGameInstance(roomId, token, auth)
	if err != nil {
		return "", err
	}
	uc.mutex.RLock()
	instance, ok := uc.gameInstances[roomId]
	uc.mutex.RUnlock()
	if !ok {
		return "", v1.ErrorAccessDenied("连接失败，游戏实例不存在")
	}
	// TODO stun server config
	_, sdpOffer, err := instance.NewConnection(auth.UserId, "stun:43.138.153.172:3478")
	if err != nil {
		uc.logger.Error("创建连接出错:", err)
		return "", v1.ErrorServiceError("连接失败，创建连接出错")
	}
	return sdpOffer, nil
}

func (uc *GameServerUseCase) SdpAnswer(_ context.Context, roomId int64, token string, auth *MemberAuthInfo, sdpAnswer string) error {
	instance, err := uc.getGameInstance(roomId, token, auth)
	if err != nil {
		return err
	}
	uc.mutex.RLock()
	instance, ok := uc.gameInstances[roomId]
	uc.mutex.RUnlock()
	if !ok {
		return v1.ErrorAccessDenied("连接失败，游戏实例不存在")
	}
	instance.mutex.RLock()
	conn, ok := instance.connections[auth.UserId]
	instance.mutex.RUnlock()
	if !ok {
		uc.logger.Error("sdpAnswer无法找到用户连接:", auth.UserId, roomId)
		return v1.ErrorServiceError("连接失败")
	}
	sdp := webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  sdpAnswer,
	}
	if err := conn.pc.SetRemoteDescription(sdp); err != nil {
		uc.logger.Error("sdpAnswer setRemoteDescription错误:", roomId, auth.UserId, err)
		return v1.ErrorServiceError("连接失败")
	}
	return nil
}

func (uc *GameServerUseCase) AddICECandidate(_ context.Context, roomId int64, token string, auth *MemberAuthInfo,
	candidate string) error {
	instance, err := uc.getGameInstance(roomId, token, auth)
	if err != nil {
		return err
	}
	instance.mutex.RLock()
	conn, ok := instance.connections[auth.UserId]
	instance.mutex.RUnlock()
	if !ok {
		uc.logger.Error("sdpAnswer无法找到用户连接:", auth.UserId, roomId)
		return v1.ErrorServiceError("连接失败")
	}
	candidateInit := webrtc.ICECandidateInit{}
	err = json.Unmarshal([]byte(candidate), &candidateInit)
	if err != nil {
		uc.logger.Error("addICE解析json错误:", err)
		return v1.ErrorServiceError("连接失败")
	}
	err = conn.pc.AddICECandidate(candidateInit)
	if err != nil {
		uc.logger.Error("addICE错误:", err)
		return v1.ErrorServiceError("addICE错误:", auth.UserId, roomId, err)
	}
	return nil
}

func (uc *GameServerUseCase) getGameInstance(roomId int64, token string, auth *MemberAuthInfo) (*GameInstance, error) {
	authInfo, err := uc.memberAuthRepo.GetAuthInfo(token, roomId)
	if err != nil {
		uc.logger.Error("获取授权错误:", err)
		return nil, v1.ErrorServiceError("验证token错误")
	}
	if authInfo == nil {
		return nil, v1.ErrorServiceError("验证token错误")
	}
	if !authInfo.equals(auth) {
		return nil, v1.ErrorAccessDenied("连接失败，token无效")
	}
	uc.mutex.RLock()
	instance, ok := uc.gameInstances[roomId]
	uc.mutex.RUnlock()
	if !ok {
		return nil, v1.ErrorAccessDenied("游戏实例不存在")
	}
	return instance, nil
}

func (uc *GameServerUseCase) GetLocalICECandidate(_ context.Context, roomId int64, token string, auth *MemberAuthInfo) ([]string, error) {
	instance, err := uc.getGameInstance(roomId, token, auth)
	if err != nil {
		return nil, err
	}
	instance.mutex.RLock()
	conn, ok := instance.connections[auth.UserId]
	instance.mutex.RUnlock()
	if !ok {
		uc.logger.Error("无法找到用户连接:", auth.UserId, roomId)
		return nil, v1.ErrorServiceError("连接失败")
	}
	return conn.getLocalICECandidates(), nil
}

func (a *MemberAuthInfo) equals(other *MemberAuthInfo) bool {
	return a.UserId == other.UserId && a.AppId == other.AppId && a.Ip == other.Ip
}

func (uc *GameServerUseCase) Restart(ctx context.Context, params RestartParams) error {
	uc.mutex.RLock()
	instance, ok := uc.gameInstances[params.RoomId]
	if !ok {
		uc.mutex.RUnlock()
		return v1.ErrorServiceError("游戏实例不存在")
	}
	uc.mutex.RUnlock()
	instance.mutex.RLock()
	defer instance.mutex.RUnlock()
	_, ok = instance.connections[params.UserId]
	if !ok {
		uc.logger.Error("无法找到用户连接:", params.UserId, params.RoomId)
		return v1.ErrorServiceError("用户未连接")
	}
	data, err := uc.gameFileRepo.GetGameData(ctx, params.GameUrl)
	if err != nil {
		return err
	}
	return instance.RestartEmulator(params.GameName, data, params.EmulatorType, params.EmulatorId, params.GameId)
}

func (uc *GameServerUseCase) createGameInstanceSession(gameInstance *GameInstance) (string, error) {
	sessionKey, _, err := uc.consul.Session().Create(&api.SessionEntry{
		TTL:      "10s",
		Behavior: "delete",
	}, nil)
	if err != nil {
		return "", err
	}
	go func() {
		e := uc.consul.Session().RenewPeriodic("10s", sessionKey, &api.WriteOptions{}, gameInstance.doneChan)
		if e != nil {
			uc.logger.Error("session renew error:", e)
		}
	}()
	return sessionKey, nil
}
