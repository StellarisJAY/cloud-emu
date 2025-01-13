package biz

import (
	"github.com/StellrisJAY/cloud-emu/gamesrv/internal/biz/game"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewGameServerUseCase, game.NewConnectionFactory)
