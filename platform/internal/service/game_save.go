package service

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"time"
)

type GameSaveService struct {
	v1.UnimplementedGameSaveServer
	gameSaveUC *biz.GameSaveUseCase
}

func NewGameSaveService(gameSaveUC *biz.GameSaveUseCase) v1.GameSaveServer {
	return &GameSaveService{
		gameSaveUC: gameSaveUC,
	}
}

func (g *GameSaveService) ListGameSave(ctx context.Context, request *v1.ListGameSaveRequest) (*v1.ListGameSaveResponse, error) {
	p := common.Pagination{
		Page:     request.Page,
		PageSize: request.PageSize,
	}
	query := biz.GameSaveQuery{RoomId: request.RoomId, EmulatorId: request.EmulatorId, GameId: request.GameId}
	saves, err := g.gameSaveUC.List(ctx, query, &p)
	if err != nil {
		e := errors.FromError(err)
		return &v1.ListGameSaveResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	result := make([]*v1.GameSaveDto, len(saves))
	for i, save := range saves {
		result[i] = &v1.GameSaveDto{
			SaveId:       save.SaveId,
			RoomId:       save.RoomId,
			EmulatorId:   save.EmulatorId,
			GameId:       save.GameId,
			RoomName:     save.RoomName,
			EmulatorName: save.EmulatorName,
			GameName:     save.GameName,
			AddTime:      save.AddTime.Format(time.DateTime),
		}
	}
	return &v1.ListGameSaveResponse{
		Code:    200,
		Message: "查询成功",
		Data:    result,
		Total:   p.Total,
	}, nil
}

func (g *GameSaveService) DeleteGameSave(ctx context.Context, request *v1.DeleteGameSaveRequest) (*v1.DeleteGameSaveResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	err := g.gameSaveUC.Delete(ctx, request.SaveId, claims.UserId)
	if err != nil {
		e := errors.FromError(err)
		return &v1.DeleteGameSaveResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.DeleteGameSaveResponse{Code: 200, Message: "删除成功"}, nil
}

func (g *GameSaveService) LoadSave(ctx context.Context, request *v1.LoadSaveRequest) (*v1.LoadSaveResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GameSaveService) SaveGame(ctx context.Context, request *v1.SaveGameRequest) (*v1.SaveGameResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	err := g.gameSaveUC.SaveGame(ctx, request.RoomId, claims.UserId)
	if err != nil {
		e := errors.FromError(err)
		return &v1.SaveGameResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.SaveGameResponse{Code: 200, Message: "保存成功"}, nil
}
