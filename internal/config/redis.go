package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"os"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		logrus.Error("Failed to connect to Redis: %v\n", err)
	} else {
		logrus.Info("Successfully connected to Redis")
	}
}
