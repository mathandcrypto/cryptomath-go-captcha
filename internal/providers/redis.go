package providers

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	redisConfig "github.com/mathandcrypto/cryptomath-go-captcha/configs/redis"
)

func NewRedisProvider(ctx context.Context, config *redisConfig.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.Address(),
		DB:   0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect redis: %w", err)
	}

	return rdb, nil
}
