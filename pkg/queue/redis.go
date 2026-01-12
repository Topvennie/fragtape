package queue

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisQueue[T any] struct {
	c      *redis.Client
	enc    Codec[T]
	key    string
	tmpKey string
}

// Interface compliance
var _ Queue[any] = (*redisQueue[any])(nil)

// Interface compliance
var _ ReliableQueue[any] = (*redisQueue[any])(nil)

type RedisCfg[T any] struct {
	Enc Codec[T]
	URL string
	Key string
	// TmpKey is the key of the temporary queue
	// If one is given then dequeue will move an item from the main queue to the temporary queue
	// Use the complete function to remove it from the temporary queue
	TmpKey string
}

func NewRedis[T any](cfg RedisCfg[T]) (Queue[T], error) {
	cfg.TmpKey = ""

	return newRedis(cfg)
}

func NewRedisReliable[T any](cfg RedisCfg[T]) (ReliableQueue[T], error) {
	if cfg.TmpKey == "" {
		return nil, errors.New("no tmp key set")
	}

	return newRedis(cfg)
}

func newRedis[T any](cfg RedisCfg[T]) (ReliableQueue[T], error) {
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
		c:      client,
		enc:    cfg.Enc,
		key:    cfg.Key,
		tmpKey: cfg.TmpKey,
	}, nil
}

func (r *redisQueue[T]) Enqueue(ctx context.Context, t T) error {
	data, err := r.enc.Marshal(t)
	if err != nil {
		return fmt.Errorf("marshal data %w", err)
	}

	return r.c.LPush(ctx, r.key, data).Err()
}

func (r *redisQueue[T]) Dequeue(ctx context.Context, opts ...DequeueOption) (T, Receipt, error) {
	if len(opts) > 0 {
		return r.dequeueBlocking(ctx, opts[0].timeout)
	}

	return r.dequeue(ctx)
}

func (r *redisQueue[T]) Complete(ctx context.Context, receipt Receipt) error {
	if r.tmpKey == "" {
		return ErrNoop
	}

	if err := r.c.LRem(ctx, r.tmpKey, 1, receipt.Raw).Err(); err != nil {
		return fmt.Errorf("mark item as complete %w", err)
	}

	return nil
}

func (r *redisQueue[T]) RequeueAll(ctx context.Context) error {
	if r.tmpKey == "" {
		return ErrNoop
	}

	var err error
	for err == nil {
		err = r.c.RPopLPush(ctx, r.tmpKey, r.key).Err()
	}

	if !errors.Is(err, redis.Nil) {
		return fmt.Errorf("rpoplpush %w", err)
	}

	return nil
}

func (r *redisQueue[T]) dequeue(ctx context.Context) (T, Receipt, error) {
	var t T
	var receipt Receipt

	var data string
	var err error

	if r.tmpKey == "" {
		// Normal queue
		data, err = r.c.RPop(ctx, r.key).Result()
	} else {
		// Reliable queue
		data, err = r.c.RPopLPush(ctx, r.key, r.tmpKey).Result()
	}
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return t, receipt, ErrEmpty
		}
		return t, receipt, fmt.Errorf("getting queue data %w", err)
	}

	t, err = r.enc.Unmarshal([]byte(data))
	if err != nil {
		return t, receipt, fmt.Errorf("unmarshal data %w", err)
	}

	return t, Receipt{Raw: data}, nil
}

func (r *redisQueue[T]) dequeueBlocking(ctx context.Context, timeout time.Duration) (T, Receipt, error) {
	var t T
	var receipt Receipt

	var data string
	var err error

	if r.tmpKey == "" {
		// Normal queue
		datas, errData := r.c.BRPop(ctx, timeout, r.key).Result()
		if len(datas) != 2 {
			return t, receipt, fmt.Errorf("unexpected brpop response length: %d", len(datas))
		}
		data = datas[1]
		err = errData
	} else {
		// Reliable queue
		data, err = r.c.BRPopLPush(ctx, r.key, r.tmpKey, timeout).Result()
	}
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return t, receipt, ErrEmpty
		}
		return t, receipt, fmt.Errorf("getting queue data %w", err)
	}

	t, err = r.enc.Unmarshal([]byte(data))
	if err != nil {
		return t, receipt, fmt.Errorf("unmarshal data %w", err)
	}

	return t, Receipt{Raw: data}, nil
}
