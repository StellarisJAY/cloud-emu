package biz

import "context"

type Emulator struct {
	EmulatorId            int64  `json:"emulatorId"`
	EmulatorName          string `json:"emulatorName"`
	Description           string `json:"description"`
	Provider              string `json:"provider"`
	SupportSave           bool   `json:"supportSave"`
	SupportGraphicSetting bool   `json:"supportGraphicSetting"`
	EmulatorType          string `json:"emulatorType"`
}

type EmulatorUseCase struct {
	repo EmulatorRepo
}

type EmulatorRepo interface {
	ListEmulator(ctx context.Context, query EmulatorQuery) ([]*Emulator, error)
	GetById(ctx context.Context, emulatorId int64) (*Emulator, error)
	GetByType(ctx context.Context, emulatorType string) (*Emulator, error)
}

type EmulatorQuery struct {
	EmulatorName          string
	Provider              string
	SupportSave           bool
	SupportGraphicSetting bool
	EmulatorType          string
}

func NewEmulatorUseCase(repo EmulatorRepo) *EmulatorUseCase {
	return &EmulatorUseCase{repo: repo}
}

func (uc *EmulatorUseCase) ListEmulator(ctx context.Context, query EmulatorQuery) ([]*Emulator, error) {
	return uc.repo.ListEmulator(ctx, query)
}
