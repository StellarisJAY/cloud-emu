package biz

import (
	"context"
	"encoding/json"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/gamesrv/internal/biz/game"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/pion/webrtc/v3"
	"sync"
)

type GameServerUseCase struct {
	memberAuthRepo MemberAuthRepo
	gameInstances  map[int64]*game.Instance
	mutex          *sync.RWMutex
	logger         *log.Helper
	consul         *api.Client
	connFactory    *game.ConnectionFactory
}

type CreateRoomInstanceParams struct {
	RoomId         int64
	Auth           *MemberAuthInfo
	RoomInstanceId int64
	EmulatorId     int64
	GameId         int64
	GameData       []byte
	EmulatorCode   string
}

type RestartParams struct {
	RoomId       int64
	UserId       int64
	EmulatorCode string
	EmulatorType string
	GameName     string
	GameUrl      string
	EmulatorId   int64
	GameId       int64
	GameData     []byte
}

type LoadSaveParams struct {
	RoomId       int64
	UserId       int64
	EmulatorCode string
	EmulatorType string
	GameName     string
	GameUrl      string
	EmulatorId   int64
	GameId       int64
	GameData     []byte
	SaveData     []byte
}

func NewGameServerUseCase(memberAuthRepo MemberAuthRepo, logger log.Logger, consul *api.Client,
	cf *game.ConnectionFactory) *GameServerUseCase {
	return &GameServerUseCase{
		memberAuthRepo: memberAuthRepo,
		logger:         log.NewHelper(logger),
		mutex:          &sync.RWMutex{},
		gameInstances:  make(map[int64]*game.Instance),
		consul:         consul,
		connFactory:    cf,
	}
}

func (uc *GameServerUseCase) CreateRoomInstance(_ context.Context, params CreateRoomInstanceParams) (string, string, error) {
	instance, err := game.MakeGameInstance(params.RoomId, params.EmulatorId, params.GameId, params.EmulatorCode, params.GameData)
	if err != nil {
		return "", "", v1.ErrorServiceError("创建游戏实例出错")
	}
	uid, _ := uuid.NewUUID()
	token := uid.String()
	if err := uc.memberAuthRepo.StoreAuthInfo(token, params.RoomId, params.Auth); err != nil {
		return "", "", v1.ErrorServiceError("创建玩家token出错")
	}
	consumerCtx, consumerCancel := context.WithCancel(context.Background())
	go instance.MessageHandler(consumerCtx)
	go instance.ListenAudioChan(consumerCtx, true)
	go instance.ListenAudioChan(consumerCtx, false)
	instance.Cancel = consumerCancel
	uc.mutex.Lock()
	defer uc.mutex.Unlock()
	uc.gameInstances[params.RoomId] = instance
	// 创建consul session
	sessionKey, err := uc.createGameInstanceSession(instance)
	return token, sessionKey, err
}

func (uc *GameServerUseCase) GetRoomInstanceToken(_ context.Context, roomId int64, auth *MemberAuthInfo) (string, error) {
	uc.mutex.RLock()
	_, ok := uc.gameInstances[roomId]
	uc.mutex.RUnlock()
	if !ok {
		return "", v1.ErrorNotFound("连接失败，游戏实例不存在")
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
	_, sdpOffer, err := uc.connFactory.NewConnection(auth.UserId, instance)
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
	conn, ok := instance.GetConnection(auth.UserId)
	if !ok {
		uc.logger.Error("sdpAnswer无法找到用户连接:", auth.UserId, roomId)
		return v1.ErrorServiceError("连接失败")
	}
	sdp := webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  sdpAnswer,
	}
	if err := conn.SetRemoteDescription(sdp); err != nil {
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
	conn, ok := instance.GetConnection(auth.UserId)
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
	err = conn.AddICECandidate(candidateInit)
	if err != nil {
		uc.logger.Error("addICE错误:", err)
		return v1.ErrorServiceError("addICE错误:", auth.UserId, roomId, err)
	}
	return nil
}

func (uc *GameServerUseCase) getGameInstance(roomId int64, token string, auth *MemberAuthInfo) (*game.Instance, error) {
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
	conn, ok := instance.GetConnection(auth.UserId)
	if !ok {
		uc.logger.Error("无法找到用户连接:", auth.UserId, roomId)
		return nil, v1.ErrorServiceError("连接失败")
	}
	return conn.GetLocalICECandidates(), nil
}

func (a *MemberAuthInfo) equals(other *MemberAuthInfo) bool {
	return a.UserId == other.UserId && a.AppId == other.AppId && a.Ip == other.Ip
}

func (uc *GameServerUseCase) Restart(_ context.Context, params RestartParams) error {
	uc.mutex.RLock()
	instance, ok := uc.gameInstances[params.RoomId]
	if !ok {
		uc.mutex.RUnlock()
		return v1.ErrorServiceError("游戏实例不存在")
	}
	uc.mutex.RUnlock()
	_, ok = instance.GetConnection(params.UserId)
	if !ok {
		uc.logger.Error("无法找到用户连接:", params.UserId, params.RoomId)
		return v1.ErrorServiceError("用户未连接")
	}
	return instance.RestartEmulator(params.GameName, params.GameData, params.EmulatorCode, params.EmulatorType, params.EmulatorId, params.GameId)
}

func (uc *GameServerUseCase) createGameInstanceSession(gameInstance *game.Instance) (string, error) {
	sessionKey, _, err := uc.consul.Session().Create(&api.SessionEntry{
		TTL:      "10s",
		Behavior: "delete",
	}, nil)
	if err != nil {
		return "", err
	}
	go func() {
		e := uc.consul.Session().RenewPeriodic("10s", sessionKey, &api.WriteOptions{}, gameInstance.DoneChan)
		if e != nil {
			uc.logger.Error("session renew error:", e)
		}
	}()
	return sessionKey, nil
}

func (uc *GameServerUseCase) SaveGame(_ context.Context, roomId, userId int64) (int64, int64, []byte, error) {
	uc.mutex.RLock()
	instance, ok := uc.gameInstances[roomId]
	uc.mutex.RUnlock()
	if !ok {
		return 0, 0, nil, v1.ErrorServiceError("游戏实例不存在")
	}
	_, ok = instance.GetConnection(userId)
	if !ok {
		uc.logger.Error("无法找到用户连接:", userId, roomId)
		return 0, 0, nil, v1.ErrorServiceError("用户未连接")
	}
	save, err := instance.SaveGame()
	if err != nil {
		return 0, 0, nil, err
	}
	return save.EmulatorId, save.GameId, save.Data, nil
}

func (uc *GameServerUseCase) LoadSave(_ context.Context, params LoadSaveParams) error {
	uc.mutex.RLock()
	instance, ok := uc.gameInstances[params.RoomId]
	uc.mutex.RUnlock()
	if !ok {
		return v1.ErrorServiceError("游戏实例不存在")
	}
	_, ok = instance.GetConnection(params.UserId)
	if !ok {
		uc.logger.Error("无法找到用户连接:", params.UserId, params.RoomId)
		return v1.ErrorServiceError("用户未连接")
	}
	err := instance.LoadSave(params.EmulatorId, params.GameId, params.EmulatorCode, params.EmulatorType,
		params.GameName, params.GameData, params.SaveData)
	return err
}

func (uc *GameServerUseCase) ListOnlineRoomMember(_ context.Context, roomId int64) ([]int64, error) {
	uc.mutex.RLock()
	instance, ok := uc.gameInstances[roomId]
	uc.mutex.RUnlock()
	if !ok {
		return nil, v1.ErrorServiceError("游戏实例不存在")
	}
	return instance.GetOnlineUsers(), nil
}

func (uc *GameServerUseCase) GetControllerPlayers(_ context.Context, roomId int64) ([]game.ControllerPlayer, error) {
	uc.mutex.RLock()
	instance, ok := uc.gameInstances[roomId]
	uc.mutex.RUnlock()
	if !ok {
		return nil, v1.ErrorServiceError("游戏实例不存在")
	}
	return instance.GetControllerPlayers(), nil
}

func (uc *GameServerUseCase) SetControllerPlayers(_ context.Context, roomId int64, players []game.ControllerPlayer) error {
	uc.mutex.RLock()
	instance, ok := uc.gameInstances[roomId]
	uc.mutex.RUnlock()
	if !ok {
		return v1.ErrorServiceError("游戏实例不存在")
	}
	instance.SetController(players)
	return nil
}

func (uc *GameServerUseCase) GetGraphicOptions(_ context.Context, roomId int64) (*game.GraphicOptions, error) {
	uc.mutex.RLock()
	instance, ok := uc.gameInstances[roomId]
	uc.mutex.RUnlock()
	if !ok {
		return nil, v1.ErrorServiceError("游戏实例不存在")
	}
	return instance.GetGraphicOptions(), nil
}

func (uc *GameServerUseCase) SetGraphicOptions(_ context.Context, roomId int64, opts *game.GraphicOptions) error {
	uc.mutex.RLock()
	instance, ok := uc.gameInstances[roomId]
	uc.mutex.RUnlock()
	if !ok {
		return v1.ErrorServiceError("游戏实例不存在")
	}
	instance.SetGraphicOptions(opts)
	return nil
}

func (uc *GameServerUseCase) ApplyMacro(_ context.Context, roomId int64, keys []string, userId int64) error {
	uc.mutex.RLock()
	instance, ok := uc.gameInstances[roomId]
	uc.mutex.RUnlock()
	if !ok {
		return v1.ErrorServiceError("游戏实例不存在")
	}
	instance.ApplyMacro(userId, keys)
	return nil
}

func (uc *GameServerUseCase) Shutdown(_ context.Context, roomId int64) error {
	uc.mutex.RLock()
	instance, ok := uc.gameInstances[roomId]
	uc.mutex.RUnlock()
	if !ok {
		return v1.ErrorServiceError("游戏实例不存在")
	}
	instance.Shutdown()
	uc.mutex.Lock()
	delete(uc.gameInstances, roomId)
	uc.mutex.Unlock()
	return nil
}

func (uc *GameServerUseCase) GetEmulatorSpeed(_ context.Context, roomId int64) (float64, error) {
	uc.mutex.RLock()
	instance, ok := uc.gameInstances[roomId]
	uc.mutex.RUnlock()
	if !ok {
		return 0, v1.ErrorServiceError("游戏实例不存在")
	}
	boost := instance.GetEmulatorSpeed()
	return boost, nil
}

func (uc *GameServerUseCase) SetEmulatorSpeed(_ context.Context, roomId int64, boost float64) (float64, error) {
	uc.mutex.RLock()
	instance, ok := uc.gameInstances[roomId]
	uc.mutex.RUnlock()
	if !ok {
		return 0, v1.ErrorServiceError("游戏实例不存在")
	}
	boost = instance.SetEmulatorSpeed(boost)
	return boost, nil
}
