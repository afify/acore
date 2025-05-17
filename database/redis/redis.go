package redis

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"

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

		if _, err := rdb.Ping(ctx).Result(); err != nil {
			slog.Error("InitRedis: could not connect to Redis", "error", err)
		} else {
			slog.Info("InitRedis: connected to Redis")
		}
	})
}

func SetRedis(key, value string, expiration time.Duration) error {
	err := rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		slog.Error("SetRedis: failed to set key", "key", key, "error", err)
	}
	return err
}

func GetRedis(key string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		slog.Error("GetRedis: failed to get key", "key", key, "error", err)
	}
	return val, err
}

func DeleteRedis(key string) error {
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		slog.Error("DeleteRedis: failed to delete key", "key", key, "error", err)
	}
	return err
}
