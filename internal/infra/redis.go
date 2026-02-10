package infra

import (
	"context"
	"errors"
	"time"

	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis(cfg *config.Config) (*redis.Client, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}
	if cfg.RedisAddr == "" {
		return nil, errors.New("redis addr is empty")
	}
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	redisClient = client
	return client, nil
}

func GetRedis() *redis.Client {
	return redisClient
}
