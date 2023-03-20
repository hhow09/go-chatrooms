package lib

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func CreateRedisClient() *redis.Client {
	redis_port := os.Getenv("REIDS_PORT")
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost:%s", redis_port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := rdb.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}

	return rdb
}
