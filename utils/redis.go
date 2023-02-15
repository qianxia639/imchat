package utils

import (
	"IMChat/utils/config"
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

func InitRedis(conf config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Address,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatal("cannot connect redis: ", err)
	}

	return client
}
