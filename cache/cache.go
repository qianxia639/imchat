package cache

import (
	"context"
	"time"
)

// ttl的默认时长为 30m
var TTL = 30 * time.Minute

type Cache interface {
	SetCache(ctx context.Context, key string, val any) error

	GetCache(ctx context.Context, key string, val any) error

	SetTtlCache(ctx context.Context, key string, val any, ttl time.Duration) error
}
