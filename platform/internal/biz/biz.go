package biz

import (
	"context"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewUserUseCase, NewRoomUseCase, NewRoomInstanceUseCase, NewRoomMemberUseCase,
	NewNotificationUseCase, NewEmulatorUseCase, NewEmulatorGameUseCase, NewButtonLayoutUseCase, NewKeyboardBindingUseCase,
	NewGameSaveUseCase, NewMacroUseCase)

type Transaction interface {
	Tx(ctx context.Context, fn func(context.Context) error) error
}
