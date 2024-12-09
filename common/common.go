package common

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewRedsync)
