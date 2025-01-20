package biz

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/bwmarrin/snowflake"
	"github.com/go-kratos/kratos/v2/log"
)

type Macro struct {
	MacroId      int64
	MacroName    string
	EmulatorType string
	AddUser      int64
	KeyCodes     []string
	ShortcutKey  string
}

type MacroQuery struct {
	EmulatorType string
	AddUser      int64
}

type MacroRepo interface {
	CreateMacro(ctx context.Context, macro *Macro) error
	GetMacro(ctx context.Context, macroId int64) (*Macro, error)
	ListMacros(ctx context.Context, query MacroQuery) ([]*Macro, error)
	ListMacrosPage(ctx context.Context, query MacroQuery, p *common.Pagination) ([]*Macro, error)
	DeleteMacro(ctx context.Context, macroId int64) error
}

type MacroUseCase struct {
	macroRepo        MacroRepo
	emulatorRepo     EmulatorRepo
	roomInstanceRepo RoomInstanceRepo
	memberRepo       RoomMemberRepo
	gameServerRepo   GameServerRepo
	snowflake        *snowflake.Node
	logger           *log.Helper
}

func NewMacroUseCase(macroRepo MacroRepo, emulatorRepo EmulatorRepo, roomInstanceRepo RoomInstanceRepo,
	memberRepo RoomMemberRepo, gameServerRepo GameServerRepo, snowflake *snowflake.Node, logger log.Logger) *MacroUseCase {
	return &MacroUseCase{
		macroRepo:        macroRepo,
		emulatorRepo:     emulatorRepo,
		roomInstanceRepo: roomInstanceRepo,
		memberRepo:       memberRepo,
		gameServerRepo:   gameServerRepo,
		snowflake:        snowflake,
		logger:           log.NewHelper(log.With(logger, "module", "biz/macro")),
	}
}

func (m *MacroUseCase) CreateMacro(ctx context.Context, macro *Macro) error {
	macro.MacroId = m.snowflake.Generate().Int64()
	return m.macroRepo.CreateMacro(ctx, macro)
}

func (m *MacroUseCase) GetMacro(ctx context.Context, macroId int64) (*Macro, error) {
	return m.macroRepo.GetMacro(ctx, macroId)
}

func (m *MacroUseCase) ListMacros(ctx context.Context, query MacroQuery) ([]*Macro, error) {
	return m.macroRepo.ListMacros(ctx, query)
}

func (m *MacroUseCase) DeleteMacro(ctx context.Context, macroId int64) error {
	return m.macroRepo.DeleteMacro(ctx, macroId)
}

func (m *MacroUseCase) ApplyMacro(ctx context.Context, macroId, roomId, userId int64) error {
	member, _ := m.memberRepo.GetByRoomAndUser(ctx, roomId, userId)
	if member == nil {
		return v1.ErrorAccessDenied("没有操作权限")
	}
	instance, _ := m.roomInstanceRepo.GetRoomInstance(ctx, roomId)
	if instance == nil {
		return v1.ErrorAccessDenied("房间实例不存在")
	}
	macro, _ := m.macroRepo.GetMacro(ctx, macroId)
	if macro == nil {
		return v1.ErrorNotFound("宏不存在")
	}
	return m.gameServerRepo.ApplyMacro(ctx, instance, macro, userId)
}
