package database

import (
	"context"
	"fmt"
	"log"

	"github.com/gulmix/Social-Network/internal/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis(cfg *config.Config) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if _, err := RedisClient.Ping(context.Background()).Result(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Successfully connected to Redis")
	return nil
}

func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
