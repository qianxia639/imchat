package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (c *RedisCache) SetCache(ctx context.Context, key string, val any) error {
	return c.client.Set(ctx, key, val, TTL).Err()
}

func (c *RedisCache) GetCache(ctx context.Context, key string, val any) error {
	return c.client.Get(ctx, key).Scan(val)
}

func (c *RedisCache) SetTtlCache(ctx context.Context, key string, val any, ttl time.Duration) error {
	return c.client.Set(ctx, key, val, ttl).Err()
}
