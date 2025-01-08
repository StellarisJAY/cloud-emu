package service

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
	"time"
)

type EmulatorService struct {
	v1.UnimplementedEmulatorServer
	uc             *biz.EmulatorUseCase
	emulatorGameUC *biz.EmulatorGameUseCase
}

func NewEmulatorService(uc *biz.EmulatorUseCase, emulatorGameUC *biz.EmulatorGameUseCase) v1.EmulatorServer {
	return &EmulatorService{uc: uc, emulatorGameUC: emulatorGameUC}
}

func (e *EmulatorService) ListEmulator(ctx context.Context, request *v1.ListEmulatorRequest) (*v1.ListEmulatorResponse, error) {
	emulators, err := e.uc.ListEmulator(ctx, biz.EmulatorQuery{
		EmulatorName:          request.EmulatorName,
		Provider:              request.Provider,
		SupportSave:           request.SupportSave,
		SupportGraphicSetting: request.SupportGraphicSetting,
	})
	if err != nil {
		e := errors.FromError(err)
		return &v1.ListEmulatorResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	emulatorList := make([]*v1.EmulatorDto, len(emulators))
	for i, emulator := range emulators {
		emulatorList[i] = &v1.EmulatorDto{
			EmulatorId:            emulator.EmulatorId,
			EmulatorName:          emulator.EmulatorName,
			Description:           emulator.Description,
			Provider:              emulator.Provider,
			SupportSave:           emulator.SupportSave,
			SupportGraphicSetting: emulator.SupportGraphicSetting,
		}
	}
	return &v1.ListEmulatorResponse{
		Code:    200,
		Message: "查询成功",
		Data:    emulatorList,
	}, nil
}

func (e *EmulatorService) ListGame(ctx context.Context, request *v1.ListGameRequest) (*v1.ListGameResponse, error) {
	p := &common.Pagination{Page: request.Page, PageSize: request.PageSize}
	games, err := e.emulatorGameUC.ListGame(ctx, biz.EmulatorGameQuery{
		EmulatorId: request.EmulatorId,
		GameName:   request.GameName,
	}, p)
	if err != nil {
		e := errors.FromError(err)
		return &v1.ListGameResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	result := make([]*v1.GameDto, len(games))
	for i, game := range games {
		result[i] = &v1.GameDto{
			EmulatorId:   game.EmulatorId,
			GameName:     game.GameName,
			GameId:       game.GameId,
			Size:         game.Size,
			CustomData:   game.CustomData,
			AddTime:      game.AddTime.Format(time.DateTime),
			EmulatorName: game.EmulatorName,
			EmulatorType: game.EmulatorType,
		}
	}
	return &v1.ListGameResponse{
		Code:    200,
		Message: "查询成功",
		Data:    result,
		Total:   p.Total,
	}, nil
}
