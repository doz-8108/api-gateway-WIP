package db

import (
	"context"
	"fmt"
	"os"

	"github.com/doz-8108/api-gateway/config"
	"github.com/redis/go-redis/v9"
)

func CreateRedisConnection(envVars config.EnvVars) *redis.Client {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     envVars.REDIS_ADDR,
		Password: envVars.REDIS_PASSWORD,
		DB:       envVars.REDIS_DB,
	})

	err := client.Conn().Ping(ctx).Err()
	if err != nil {
		fmt.Println("Failed to connect to Redis")
		os.Exit(0)
	}

	return client
}
