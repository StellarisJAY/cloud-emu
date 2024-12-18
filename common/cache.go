package common

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache[T any] struct {
	redis *redis.Client
}

const keyPrefix = "cloud-emu"

func NewCache[T any](redis *redis.Client) *Cache[T] {
	return &Cache[T]{
		redis: redis,
	}
}

func (c *Cache[T]) Set(ctx context.Context, key string, value *T, expire time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, keyPrefix+key, data, expire).Err()
}

func (c *Cache[T]) Get(ctx context.Context, key string) (*T, error) {
	cmd := c.redis.Get(ctx, keyPrefix+key)
	if errors.Is(cmd.Err(), redis.Nil) {
		return nil, nil
	} else if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	value, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}
	res := new(T)
	if err := json.Unmarshal(value, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Cache[T]) Del(ctx context.Context, key string) error {
	return c.redis.Del(ctx, keyPrefix+key).Err()
}
