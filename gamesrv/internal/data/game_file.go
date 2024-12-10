package data

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/gamesrv/internal/biz"
)

type GameFileRepo struct {
	data *Data
}

func NewGameFileRepo(data *Data) biz.GameFileRepo {
	return &GameFileRepo{
		data: data,
	}
}

func (g *GameFileRepo) GetGameData(ctx context.Context, game string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GameFileRepo) GetSavedGame(ctx context.Context, id int64) (*biz.GameSave, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GameFileRepo) SaveGame(ctx context.Context, save *biz.GameSave) error {
	//TODO implement me
	panic("implement me")
}

func (g *GameFileRepo) ListSaves(ctx context.Context, roomId int64, page, pageSize int32) ([]*biz.GameSave, int32, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GameFileRepo) DeleteSave(ctx context.Context, saveId int64) error {
	//TODO implement me
	panic("implement me")
}

func (g *GameFileRepo) GetExitSave(ctx context.Context, roomId int64) (*biz.GameSave, error) {
	//TODO implement me
	panic("implement me")
}
