package data

import (
	"github.com/StellrisJAY/cloud-emu/platform/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewRedisClient, NewUserRepo, NewRoomRepo, NewRoomInstanceRepo, NewAuthRepo,
	NewGameServerRepo, NewRoomMemberRepo)

// Data .
type Data struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewRedisClient(c *conf.Data) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
	})
}

// NewData .
func NewData(c *conf.Data, redis *redis.Client, logger log.Logger) (*Data, func(), error) {
	d := &Data{}
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		if d.redis != nil {
			_ = d.redis.Close()
		}
	}
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		Logger: logger2.Default.LogMode(logger2.Info),
	})

	if err != nil {
		return nil, cleanup, err
	}
	d.db = db
	d.redis = redis
	return d, cleanup, nil
}
