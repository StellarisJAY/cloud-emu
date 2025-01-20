package service

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
)

type MacroService struct {
	v1.UnimplementedMacroServer
	uc *biz.MacroUseCase
}

func NewMacroService(uc *biz.MacroUseCase) v1.MacroServer {
	return &MacroService{uc: uc}
}

func (m *MacroService) ListMacros(ctx context.Context, request *v1.ListMacrosRequest) (*v1.ListMacrosResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	macros, err := m.uc.ListMacros(ctx, biz.MacroQuery{
		EmulatorType: request.EmulatorType,
		AddUser:      claims.UserId,
	})
	if err != nil {
		e := errors.FromError(err)
		return &v1.ListMacrosResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	result := make([]*v1.MacroDto, len(macros))
	for i, macro := range macros {
		result[i] = &v1.MacroDto{
			MacroId:      macro.MacroId,
			MacroName:    macro.MacroName,
			EmulatorType: macro.EmulatorType,
			AddUser:      macro.AddUser,
			KeyCodes:     macro.KeyCodes,
			ShortcutKey:  macro.ShortcutKey,
		}
	}
	return &v1.ListMacrosResponse{
		Code:    200,
		Message: "操作成功",
		Data:    result,
	}, nil
}

func (m *MacroService) CreateMacro(ctx context.Context, request *v1.CreateMacroRequest) (*v1.CreateMacroResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	err := m.uc.CreateMacro(ctx, &biz.Macro{
		MacroName:    request.MacroName,
		EmulatorType: request.EmulatorType,
		AddUser:      claims.UserId,
		KeyCodes:     request.KeyCodes,
		ShortcutKey:  request.ShortcutKey,
	})
	if err != nil {
		e := errors.FromError(err)
		return &v1.CreateMacroResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.CreateMacroResponse{
		Code:    200,
		Message: "操作成功",
	}, nil
}

func (m *MacroService) DeleteMacro(ctx context.Context, request *v1.DeleteMacroRequest) (*v1.DeleteMacroResponse, error) {
	err := m.uc.DeleteMacro(ctx, request.MacroId)
	if err != nil {
		e := errors.FromError(err)
		return &v1.DeleteMacroResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.DeleteMacroResponse{
		Code:    200,
		Message: "操作成功",
	}, nil
}

func (m *MacroService) ApplyMacro(ctx context.Context, request *v1.ApplyMacroRequest) (*v1.ApplyMacroResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	err := m.uc.ApplyMacro(ctx, request.MacroId, request.RoomId, claims.UserId)
	if err != nil {
		e := errors.FromError(err)
		return &v1.ApplyMacroResponse{
			Code:    e.Code,
			Message: e.Message,
		}, nil
	}
	return &v1.ApplyMacroResponse{
		Code:    200,
		Message: "操作成功",
	}, nil
}
