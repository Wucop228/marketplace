package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedis(redisAddr string, redisPassword string, redisDB int) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})
	return &Redis{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *Redis) Set(key string, value string, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *Redis) Del(key string) error {
	return r.client.Del(r.ctx, key).Err()
}
