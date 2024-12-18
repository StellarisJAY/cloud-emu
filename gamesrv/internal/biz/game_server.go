package biz

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/google/uuid"
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
}

type CreateRoomInstanceParams struct {
	RoomId     int64
	EmulatorId int64
	GameFile   string
	Auth       *MemberAuthInfo
}

func NewGameServerUseCase(gameFileRepo GameFileRepo, memberAuthRepo MemberAuthRepo) *GameServerUseCase {
	return &GameServerUseCase{
		gameFileRepo:   gameFileRepo,
		memberAuthRepo: memberAuthRepo,
	}
}

func (uc *GameServerUseCase) CreateRoomInstance(ctx context.Context, params CreateRoomInstanceParams) error {
	instance, err := makeGameInstance(params.RoomId)
	if err != nil {
		return v1.ErrorServiceError("创建游戏实例出错")
	}
	uid, _ := uuid.NewUUID()
	if err := uc.memberAuthRepo.StoreAuthInfo(uid.String(), params.RoomId, params.Auth); err != nil {
		return v1.ErrorServiceError("创建玩家token出错")
	}
	uc.mutex.Lock()
	defer uc.mutex.Unlock()
	uc.gameInstances[params.RoomId] = instance
	return nil
}

func (uc *GameServerUseCase) GetRoomInstanceToken(ctx context.Context, roomId int64, auth *MemberAuthInfo) (string, error) {
	panic("not implemented")
}
