package db

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client     *redis.Client
	onceRedis  sync.Once
	defaultCtx = context.Background()
)

// Initialize sets up the Redis client only once
func InitRedis() {
	onceRedis.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr: os.Getenv("REDIS_ADDR"),
			DB:   0,
		})

		if err := ping(); err != nil {
			panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
		}
	})
}

func ping() error {
	ctx, cancel := context.WithTimeout(defaultCtx, 5*time.Second)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	return err
}

// GetClient returns the Redis client instance
func GetClientRedis() *redis.Client {
	if client == nil {
		panic("Redis client not initialized. Call redisclient.Initialize() first.")
	}
	return client
}
