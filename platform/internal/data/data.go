package data

import (
	"github.com/StellrisJAY/cloud-emu/platform/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewRoomRepo, NewRoomInstanceRepo, NewAuthRepo)

// Data .
type Data struct {
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	d := &Data{}
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		Logger: logger2.Default.LogMode(logger2.Info),
	})
	if err != nil {
		return nil, cleanup, err
	}
	d.db = db
	return d, cleanup, nil
}
