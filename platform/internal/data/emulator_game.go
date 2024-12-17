package data

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"time"
)

type EmulatorGameEntity struct {
	GameId     int64
	EmulatorId int64
	GameName   string
	Size       int32
	AddTime    time.Time
	CustomData string
	AddUser    int64
	Url        string
}

type EmulatorGameRepo struct {
	data *Data
}

const EmulatorGameTableName = "emulator_game"

func NewEmulatorGameRepo(data *Data) biz.EmulatorGameRepo {
	return &EmulatorGameRepo{data: data}
}

func (e *EmulatorGameRepo) Create(ctx context.Context, game *biz.EmulatorGame) error {
	return e.data.DB(ctx).Table(EmulatorGameTableName).WithContext(ctx).Create(convertEmulatorGameBizToEntity(game)).Error
}

func (e *EmulatorGameRepo) Upload(ctx context.Context, game *biz.EmulatorGame, data []byte) error {
	//TODO upload game to file system
	return nil
}

func (e *EmulatorGameRepo) Delete(ctx context.Context, gameId int64) error {
	//TODO delete game file
	panic("implement me")
}

func (e *EmulatorGameRepo) ListGame(ctx context.Context, query biz.EmulatorGameQuery, p *common.Pagination) ([]*biz.EmulatorGame, error) {
	var result []*biz.EmulatorGame
	d := e.data.DB(ctx).Table(EmulatorGameTableName).WithContext(ctx)
	if query.EmulatorId != 0 {
		d = d.Where("emulator_id = ?", query.EmulatorId)
	}
	if query.GameName != "" {
		d = d.Where("game_name LIKE ?", "%"+query.GameName+"%")
	}
	if err := d.Scopes(common.WithPagination(p)).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func convertEmulatorGameBizToEntity(gameBiz *biz.EmulatorGame) *EmulatorGameEntity {
	return &EmulatorGameEntity{
		GameId:     gameBiz.GameId,
		EmulatorId: gameBiz.EmulatorId,
		GameName:   gameBiz.GameName,
		Size:       gameBiz.Size,
		AddTime:    gameBiz.AddTime,
		CustomData: gameBiz.CustomData,
		AddUser:    gameBiz.AddUser,
		Url:        gameBiz.Url,
	}
}
