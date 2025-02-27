package db

import (
	"context"

	"github.com/doz-8108/api-gateway/config"
	"github.com/doz-8108/api-gateway/utils"
	"github.com/redis/go-redis/v9"
)

func CreateRedisConnection(envVars config.EnvVars, utils utils.Utils) *redis.Client {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     envVars.REDIS_ADDR,
		Password: envVars.REDIS_PASSWORD,
		DB:       envVars.REDIS_DB,
		Protocol: 2,
	})

	err := client.Conn().Ping(ctx).Err()
	if err != nil {
		utils.Logger.Fatalf("Failed to connect to Redis: %v", err)
	}

	return client
}
