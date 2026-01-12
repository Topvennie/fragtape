package queue

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisQueue[T any] struct {
	c   *redis.Client
	enc Codec[T]
	key string
}

// Interface compliance
var _ Queue[any] = (*redisQueue[any])(nil)

type RedisCfg[T any] struct {
	Enc Codec[T]
	URL string
	Key string
}

func NewRedis[T any](cfg RedisCfg[T]) (Queue[T], error) {
	options, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	if cfg.Enc == nil {
		return nil, errors.New("no encoder set")
	}

	if cfg.Key == "" {
		return nil, errors.New("no key set")
	}

	return &redisQueue[T]{
		c:   client,
		enc: cfg.Enc,
		key: cfg.Key,
	}, nil
}

func (r *redisQueue[T]) Enqueue(ctx context.Context, t T) error {
	data, err := r.enc.Marshal(t)
	if err != nil {
		return fmt.Errorf("marshal data %w", err)
	}

	return r.c.LPush(ctx, r.key, data).Err()
}

func (r *redisQueue[T]) Dequeue(ctx context.Context, opts ...DequeueOption) (T, error) {
	if len(opts) > 0 {
		return r.dequeueBlocking(ctx, opts[0].timeout)
	}

	return r.dequeue(ctx)
}

func (r *redisQueue[T]) dequeue(ctx context.Context) (T, error) {
	var t T

	data, err := r.c.RPop(ctx, r.key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return t, ErrEmpty
		}
		return t, fmt.Errorf("getting queue data %w", err)
	}

	t, err = r.enc.Unmarshal([]byte(data))
	if err != nil {
		return t, fmt.Errorf("unmarshal data %w", err)
	}

	return t, nil
}

func (r *redisQueue[T]) dequeueBlocking(ctx context.Context, timeout time.Duration) (T, error) {
	var t T

	data, err := r.c.BRPop(ctx, timeout, r.key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return t, ErrEmpty
		}
		return t, fmt.Errorf("getting queue data %w", err)
	}
	if len(data) != 2 {
		return t, fmt.Errorf("unexpected brpop response length: %d", len(data))
	}

	t, err = r.enc.Unmarshal([]byte(data[1]))
	if err != nil {
		return t, fmt.Errorf("unmarshal data %w", err)
	}

	return t, nil
}
