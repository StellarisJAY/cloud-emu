package service

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewUserService, NewRoomService, NewRoomInstanceService, NewNotificationService,
	NewRoomMemberService, NewEmulatorService, NewButtonLayoutService, NewKeyboardBindingService, NewGameSaveService,
	NewMacroService)
