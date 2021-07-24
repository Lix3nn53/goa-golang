package storage

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

//RedisCtx ...
var RedisCtx = context.Background()

// DbCache ...
type DbCache struct {
	*redis.Client
}

//InitializeCache ...
func InitializeCache() *DbCache {
	db, _ := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	rdb := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     fmt.Sprintf(os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})
	return &DbCache{rdb}
}
