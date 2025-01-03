package common

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

func NewRedsync(cli *goredislib.Client) *redsync.Redsync {
	pool := goredis.NewPool(cli)
	return redsync.New(pool)
}
