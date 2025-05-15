package database

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	rdb  *redis.Client
	once sync.Once
	ctx  = context.Background()
)

func InitRedis() {
	once.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})

		_, err := rdb.Ping(ctx).Result()
		if err != nil {
			log.Printf("Could not connect to Redis: %v", err)
		}
		log.Println("[1] Connected to Redis")
	})
}

func SetRedis(key string, value interface{}) error {
	return rdb.Set(ctx, key, value, 0).Err()
}

// GetRedis gets a value by key from Redis
func GetRedis(key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}

// DeleteRedis deletes a key from Redis
func DeleteRedis(key string) error {
	return rdb.Del(ctx, key).Err()
}
