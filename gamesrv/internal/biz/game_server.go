package biz

import (
	"context"
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

type GameFileMetadata struct {
	Name       string `json:"name"`
	Mapper     string `json:"mapper"`
	Mirroring  string `json:"mirroring"`
	Size       int32  `json:"size"`
	UploadTime int64  `json:"uploadTime"`
}

type GameInstanceStats struct {
	RoomId            int64         `json:"roomId"`
	Connections       int           `json:"connections"`
	ActiveConnections int           `json:"activeConnections"`
	Game              string        `json:"game"`
	Uptime            time.Duration `json:"uptime"`
}

type EndpointStats struct {
	EmulatorCount int32 `json:"emulatorCount"`
	CpuUsage      int32 `json:"cpuUsage"`
	MemoryUsed    int64 `json:"memoryUsed"`
	MemoryTotal   int64 `json:"memoryTotal"`
	Uptime        int64 `json:"uptime"`
}

type GraphicOptions struct {
	HighResOpen  bool `json:"highResOpen"`
	ReverseColor bool `json:"reverseColor"`
	Grayscale    bool `json:"grayscale"`
}

type SupportedEmulator struct {
	Name                  string `json:"name"`
	SupportSaving         bool   `json:"saving"`
	SupportGraphicSetting bool   `json:"supportGraphicSetting"`
	Provider              string `json:"provider"`
	Games                 int32  `json:"games"`
}

type GameFileRepo interface {
	GetGameData(ctx context.Context, game string) ([]byte, error)
	GetSavedGame(ctx context.Context, id int64) (*GameSave, error)
	SaveGame(ctx context.Context, save *GameSave) error
	ListSaves(ctx context.Context, roomId int64, page, pageSize int32) ([]*GameSave, int32, error)
	DeleteSave(ctx context.Context, saveId int64) error
	GetExitSave(ctx context.Context, roomId int64) (*GameSave, error)
}
