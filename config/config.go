package config

import (
	"context"
	redis "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg *AppConfig) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	_, err := client.Ping(context.Background()).Result() // Test connection
	if err != nil {
		panic(err)
	}

	return &RedisClient{Client: client}
}