package biz

import "context"

type Emulator struct {
	EmulatorId            int64
	EmulatorName          string
	Description           string
	Provider              string
	SupportSave           bool
	SupportGraphicSetting bool
}

type EmulatorUseCase struct {
	repo EmulatorRepo
}

type EmulatorRepo interface {
	ListEmulator(ctx context.Context, query EmulatorQuery) ([]*Emulator, error)
}

type EmulatorQuery struct {
	EmulatorName          string
	Provider              string
	SupportSave           bool
	SupportGraphicSetting bool
}

func NewEmulatorUseCase(repo EmulatorRepo) *EmulatorUseCase {
	return &EmulatorUseCase{repo: repo}
}

func (uc *EmulatorUseCase) ListEmulator(ctx context.Context, query EmulatorQuery) ([]*Emulator, error) {
	return uc.repo.ListEmulator(ctx, query)
}
