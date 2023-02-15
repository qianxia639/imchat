package utils

import (
	"IMChat/utils/config"

	"github.com/go-redis/redis/v8"
)

func InitRedis(conf config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Address,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})
}
