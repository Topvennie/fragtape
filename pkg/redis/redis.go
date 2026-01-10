// Package redis connects to the redis db
package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	C      *redis.Client
	ErrNil = redis.Nil
)

type RedisCfg struct {
	URL string
}

func New(cfg RedisCfg) error {
	options, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return err
	}

	C = redis.NewClient(options)
	ctx := context.Background()
	_, err = C.Ping(ctx).Result()
	return err
}
