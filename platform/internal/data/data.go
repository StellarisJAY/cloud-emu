package data

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/platform/internal/biz"
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
	NewGameServerRepo, NewRoomMemberRepo, NewNotificationRepo, NewUserEmailVerifyRepo, NewTransaction, NewEmulatorRepo,
	NewEmulatorGameRepo)

// Data .
type Data struct {
	db    *gorm.DB
	redis *redis.Client
}

type txKey struct{}

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

func (d *Data) DB(ctx context.Context) *gorm.DB {
	db, ok := ctx.Value(txKey{}).(*gorm.DB)
	if !ok {
		return d.db
	}
	return db
}

func (d *Data) Tx(ctx context.Context, fn func(c context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(db *gorm.DB) error {
		ctx = context.WithValue(ctx, txKey{}, db)
		return fn(ctx)
	})
}

func NewTransaction(d *Data) biz.Transaction {
	return d
}
