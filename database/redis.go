package database

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// RedisClient is the Redis client
var RedisClient *redis.Client

// ConnectRedis connects to Redis
func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   4,
	})

	err := RedisClient.Ping(context.Background()).Err()

	return err
}

// DisconnectRedis disconnects from Redis
func DisconnectRedis() error {
	if err := RedisClient.Close(); err != nil {
		return err
	}

	return nil
}
