package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr, password string, db int) *RedisClient {
	return &RedisClient{Client: redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db})}
}

func (r *RedisClient) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}
