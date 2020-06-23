package redis

import (
	"context"
	"fmt"
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
	pong, err := rdb.Ping(ctx).Result()
	fmt.Println(pong, err)
	if err != nil {
		panic(err)
	}
	return rdb
}
