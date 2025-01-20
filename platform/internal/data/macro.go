package data

import (
	"context"
	"errors"
	"github.com/StellrisJAY/cloud-emu/common"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
	"gorm.io/gorm"
	"strings"
)

type Macro struct {
	MacroId      int64
	MacroName    string
	EmulatorType string
	AddUser      int64
	KeyCodes     string
	ShortcutKey  string
}

type MacroRepo struct {
	data *Data
}

const MacroTableName = "macros"

func NewMacroRepo(data *Data) biz.MacroRepo {
	return &MacroRepo{data: data}
}

func (m *MacroRepo) CreateMacro(ctx context.Context, macro *biz.Macro) error {
	return m.data.DB(ctx).
		Table(MacroTableName).
		Create(convertMacroBizToEntity(macro)).
		WithContext(ctx).
		Error
}

func (m *MacroRepo) GetMacro(ctx context.Context, macroId int64) (*biz.Macro, error) {
	var macro *Macro
	err := m.data.DB(ctx).
		Table(MacroTableName).
		Where("macro_id = ?", macroId).
		WithContext(ctx).
		First(&macro).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return convertMacroEntityToBiz(macro), nil
}

func (m *MacroRepo) ListMacros(ctx context.Context, query biz.MacroQuery) ([]*biz.Macro, error) {
	var macros []*Macro
	err := m.data.DB(ctx).
		Table(MacroTableName).
		Where("emulator_type = ? and add_user = ?", query.EmulatorType, query.AddUser).
		WithContext(ctx).
		Scan(&macros).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	res := make([]*biz.Macro, len(macros))
	for i, macro := range macros {
		res[i] = convertMacroEntityToBiz(macro)
	}
	return res, nil
}

func (m *MacroRepo) ListMacrosPage(ctx context.Context, query biz.MacroQuery, p *common.Pagination) ([]*biz.Macro, error) {
	var macros []*Macro
	err := m.data.DB(ctx).
		Table(MacroTableName).
		Where("emulator_type = ? and add_user = ?", query.EmulatorType, query.AddUser).
		Scopes(common.WithPagination(p)).
		WithContext(ctx).
		Scan(&macros).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	res := make([]*biz.Macro, len(macros))
	for i, macro := range macros {
		res[i] = convertMacroEntityToBiz(macro)
	}
	return res, nil
}

func (m *MacroRepo) DeleteMacro(ctx context.Context, macroId int64) error {
	return m.data.DB(ctx).
		Table(MacroTableName).
		Where("macro_id =?", macroId).
		WithContext(ctx).Delete(&Macro{}).Error
}

func convertMacroEntityToBiz(macro *Macro) *biz.Macro {
	return &biz.Macro{
		MacroId:      macro.MacroId,
		MacroName:    macro.MacroName,
		EmulatorType: macro.EmulatorType,
		AddUser:      macro.AddUser,
		ShortcutKey:  macro.ShortcutKey,
		KeyCodes:     strings.Split(macro.KeyCodes, ","),
	}
}

func convertMacroBizToEntity(macro *biz.Macro) *Macro {
	return &Macro{
		MacroId:      macro.MacroId,
		MacroName:    macro.MacroName,
		EmulatorType: macro.EmulatorType,
		AddUser:      macro.AddUser,
		ShortcutKey:  macro.ShortcutKey,
		KeyCodes:     strings.Join(macro.KeyCodes, ","),
	}
}
