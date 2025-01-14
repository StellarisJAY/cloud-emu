package data

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/gamesrv/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewMemberAuthRepo)

// Data .
type Data struct {
	redis   *redis.Client
	mongodb *mongo.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	d := &Data{}
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		if d.redis != nil {
			_ = d.redis.Close()
		}
		if d.mongodb != nil {
			_ = d.mongodb.Disconnect(context.Background())
		}
	}
	d.redis = redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
	})
	opts := options.Client().SetHosts([]string{c.Mongodb.Host})
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, cleanup, err
	}
	d.mongodb = client
	return d, cleanup, nil
}

func (d *Data) getGridFSBucket(dbName string, bucketName string) (*gridfs.Bucket, error) {
	bucket, err := gridfs.NewBucket(d.mongodb.Database(dbName), options.GridFSBucket().SetName(bucketName))
	if err != nil {
		return nil, err
	}
	return bucket, nil
}
