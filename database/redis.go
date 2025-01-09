package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	RedisDB *redis.Client
	RedisContext   = context.Background()
	RedisCacheTTL = 10 * time.Second // Cache time-to-live
)

func SetupConnectRedis() {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})

	_, err := RedisDB.Ping(RedisDB.Context()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	fmt.Println("Connection Opened to Redis")
}
