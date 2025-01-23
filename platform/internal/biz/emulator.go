package biz

import (
	"context"
	v1 "github.com/StellrisJAY/cloud-emu/api/v1"
	"github.com/StellrisJAY/cloud-emu/emulator"
)

type Emulator struct {
	EmulatorId            int64  `json:"emulatorId"`
	EmulatorName          string `json:"emulatorName"`
	Description           string `json:"description"`
	Provider              string `json:"provider"`
	SupportSave           bool   `json:"supportSave"`
	SupportGraphicSetting bool   `json:"supportGraphicSetting"`
	EmulatorType          string `json:"emulatorType"`
	EmulatorCode          string `json:"emulatorCode"`
	Disabled              bool   `json:"disabled"`
}

type EmulatorUseCase struct {
	repo     EmulatorRepo
	userRepo UserRepo
}

type EmulatorRepo interface {
	ListEmulator(ctx context.Context, query EmulatorQuery) ([]*Emulator, error)
	GetById(ctx context.Context, emulatorId int64) (*Emulator, error)
	GetByType(ctx context.Context, emulatorType string) (*Emulator, error)
	Update(ctx context.Context, emulator *Emulator) error
}

type EmulatorQuery struct {
	EmulatorName          string
	Provider              string
	SupportSave           bool
	SupportGraphicSetting bool
	EmulatorType          string
	ShowDisabled          bool
}

func NewEmulatorUseCase(repo EmulatorRepo, userRepo UserRepo) *EmulatorUseCase {
	return &EmulatorUseCase{repo: repo, userRepo: userRepo}
}

func (uc *EmulatorUseCase) ListEmulator(ctx context.Context, query EmulatorQuery) ([]*Emulator, error) {
	return uc.repo.ListEmulator(ctx, query)
}

func (uc *EmulatorUseCase) ListEmulatorTypes(_ context.Context) []string {
	return emulator.GetSupportedEmulatorTypes()
}

func (uc *EmulatorUseCase) Update(ctx context.Context, emulator *Emulator, userId int64) error {
	user, _ := uc.userRepo.GetById(ctx, userId)
	if user == nil || user.Role != UserRoleAdmin {
		return v1.ErrorAccessDenied("没有修改权限")
	}
	return uc.repo.Update(ctx, emulator)
}
