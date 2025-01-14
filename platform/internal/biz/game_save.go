package biz

import (
	"context"
	"fmt"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/bwmarrin/snowflake"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type GameSave struct {
	SaveId       int64
	RoomId       int64
	EmulatorId   int64
	GameId       int64
	RoomName     string
	EmulatorName string
	EmulatorType string
	EmulatorCode string
	GameName     string
	AddTime      time.Time
	FileUrl      string
}

type GameSaveQuery struct {
	RoomId     int64
	EmulatorId int64
	GameId     int64
}

type GameSaveRepo interface {
	Create(ctx context.Context, save *GameSave) error
	Upload(ctx context.Context, save *GameSave, data []byte) error
	Download(ctx context.Context, save *GameSave) ([]byte, error)
	List(ctx context.Context, query GameSaveQuery, p *common.Pagination) ([]*GameSave, error)
	Delete(ctx context.Context, saveId int64) error
	DeleteFile(ctx context.Context, save *GameSave) error
	Get(ctx context.Context, saveId int64) (*GameSave, error)
	GetDetail(ctx context.Context, saveId int64) (*GameSave, error)
}

type GameSaveUseCase struct {
	gameSaveRepo     GameSaveRepo
	roomInstanceRepo RoomInstanceRepo
	gameServerRepo   GameServerRepo
	roomMemberRepo   RoomMemberRepo
	emulatorGameRepo EmulatorGameRepo
	emulatorRepo     EmulatorRepo
	snowflake        *snowflake.Node
	tm               Transaction
	logger           *log.Helper
}

func NewGameSaveUseCase(gameSaveRepo GameSaveRepo, roomInstanceRepo RoomInstanceRepo, gameServerRepo GameServerRepo,
	roomMemberRepo RoomMemberRepo, emulatorGameRepo EmulatorGameRepo, emulatorRepo EmulatorRepo, snowflake *snowflake.Node,
	tm Transaction, logger log.Logger) *GameSaveUseCase {
	return &GameSaveUseCase{
		gameSaveRepo:     gameSaveRepo,
		roomInstanceRepo: roomInstanceRepo,
		gameServerRepo:   gameServerRepo,
		roomMemberRepo:   roomMemberRepo,
		emulatorGameRepo: emulatorGameRepo,
		emulatorRepo:     emulatorRepo,
		snowflake:        snowflake,
		tm:               tm,
		logger:           log.NewHelper(logger),
	}
}

func (uc *GameSaveUseCase) List(ctx context.Context, query GameSaveQuery, p *common.Pagination) ([]*GameSave, error) {
	return uc.gameSaveRepo.List(ctx, query, p)
}

func (uc *GameSaveUseCase) Delete(ctx context.Context, saveId int64, userId int64) error {
	return uc.tm.Tx(ctx, func(ctx context.Context) error {
		save, _ := uc.gameSaveRepo.Get(ctx, saveId)
		if save == nil {
			return v1.ErrorNotFound("存档不存在")
		}
		rm, _ := uc.roomMemberRepo.GetByRoomAndUser(ctx, save.RoomId, userId)
		if rm == nil || rm.Role != RoomMemberRoleHost {
			return v1.ErrorAccessDenied("没有权限")
		}
		if err := uc.gameSaveRepo.Delete(ctx, saveId); err != nil {
			uc.logger.Error("删除存档sql出错:", err)
			return v1.ErrorServiceError("删除出错")
		}
		if err := uc.gameSaveRepo.DeleteFile(ctx, save); err != nil {
			uc.logger.Error("删除存档文件出错:", err)
			return v1.ErrorServiceError("删除出错")
		}
		return nil
	})
}

func (uc *GameSaveUseCase) SaveGame(ctx context.Context, roomId, userId int64) error {
	return uc.tm.Tx(ctx, func(ctx context.Context) error {
		rm, _ := uc.roomMemberRepo.GetByRoomAndUser(ctx, roomId, userId)
		if rm == nil || rm.Role != RoomMemberRoleHost {
			return v1.ErrorAccessDenied("没有权限")
		}
		ri, _ := uc.roomInstanceRepo.GetRoomInstance(ctx, roomId)
		if ri == nil {
			return v1.ErrorNotFound("房间不存在")
		}

		e, _ := uc.emulatorRepo.GetById(ctx, ri.EmulatorId)
		if e == nil || !e.SupportSave {
			return v1.ErrorServiceError("模拟器不支持存档功能")
		}

		// 游戏服务保存游戏，返回存档数据
		emulatorId, gameId, saveData, err := uc.gameServerRepo.SaveGame(ctx, ri, roomId, userId)
		if err != nil {
			uc.logger.Error("保存游戏获取存档数据错误:", err)
			return v1.ErrorServiceError("保存出错")
		}
		// 数据库保存存档信息
		save := &GameSave{
			SaveId:     uc.snowflake.Generate().Int64(),
			RoomId:     roomId,
			EmulatorId: emulatorId,
			GameId:     gameId,
			AddTime:    time.Now(),
		}
		save.FileUrl = fmt.Sprintf("mongo://cloud-emu/game_save/%d", save.SaveId)
		if err := uc.gameSaveRepo.Create(ctx, save); err != nil {
			uc.logger.Error("保存游戏sql错误:", err)
			return v1.ErrorServiceError("保存出错")
		}
		// 上传存档文件
		if err := uc.gameSaveRepo.Upload(ctx, save, saveData); err != nil {
			uc.logger.Error("保存游戏上传文件错误:", err)
			return v1.ErrorServiceError("保存出错")
		}
		return nil
	})
}

func (uc *GameSaveUseCase) LoadSave(ctx context.Context, roomId, userId int64, saveId int64) error {
	rm, _ := uc.roomMemberRepo.GetByRoomAndUser(ctx, roomId, userId)
	if rm == nil || rm.Role != RoomMemberRoleHost {
		return v1.ErrorAccessDenied("没有权限")
	}
	ri, _ := uc.roomInstanceRepo.GetRoomInstance(ctx, roomId)
	if ri == nil {
		return v1.ErrorNotFound("房间不存在")
	}
	save, err := uc.gameSaveRepo.GetDetail(ctx, saveId)
	if err != nil {
		uc.logger.Error("获取存档信息错误:", err)
		return v1.ErrorServiceError("获取存档信息错误")
	}
	if save == nil {
		return v1.ErrorNotFound("存档不存在")
	}

	params := LoadSaveParams{
		UserId:       userId,
		EmulatorId:   save.EmulatorId,
		EmulatorCode: save.EmulatorCode,
		GameId:       save.GameId,
		GameName:     save.GameName,
		EmulatorType: save.EmulatorType,
	}

	// 加载存档数据
	saveData, err := uc.gameSaveRepo.Download(ctx, save)
	if err != nil {
		uc.logger.Error("获取存档文件错误:", err)
		return v1.ErrorServiceError("读取存档错误")
	}
	params.SaveData = saveData

	// 存档模拟器与当前模拟器不同，或游戏不同，需要加载游戏文件并重启
	if save.EmulatorId != ri.EmulatorId || save.GameId != ri.GameId {
		gameData, err := uc.emulatorGameRepo.Download(ctx, &EmulatorGame{EmulatorId: save.EmulatorId, GameId: save.GameId})
		if err != nil {
			uc.logger.Error("获取游戏文件错误:", err)
			return v1.ErrorServiceError("读取存档错误")
		}
		params.GameData = gameData
	}

	if err := uc.gameServerRepo.LoadSave(ctx, ri, params); err != nil {
		uc.logger.Error("加载存档错误:", err)
		return v1.ErrorServiceError("加载存档错误")
	}
	return nil
}
