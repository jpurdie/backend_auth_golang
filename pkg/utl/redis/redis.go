package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
)

var ctx = context.Background()

func BuildRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PW"),
		DB:       0, // use default DB
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	return rdb
}
