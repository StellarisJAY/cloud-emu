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
}

type GameSaveUseCase struct {
	gameSaveRepo     GameSaveRepo
	roomInstanceRepo RoomInstanceRepo
	gameServerRepo   GameServerRepo
	roomMemberRepo   RoomMemberRepo
	snowflake        *snowflake.Node
	tm               Transaction
	logger           *log.Helper
}

func NewGameSaveUseCase(gameSaveRepo GameSaveRepo, roomInstanceRepo RoomInstanceRepo, gameServerRepo GameServerRepo,
	roomMemberRepo RoomMemberRepo, snowflake *snowflake.Node, tm Transaction, logger log.Logger) *GameSaveUseCase {
	return &GameSaveUseCase{
		gameSaveRepo:     gameSaveRepo,
		roomInstanceRepo: roomInstanceRepo,
		gameServerRepo:   gameServerRepo,
		roomMemberRepo:   roomMemberRepo,
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
