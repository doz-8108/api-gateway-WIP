package middlewares

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/storage/redis/v3"
)

func (m *Middlewares) getStorageConfig() *redis.Storage {
	redisAddrParts := strings.Split(m.EnvVars.REDIS_ADDR, ":")
	redisPort, err := strconv.Atoi(redisAddrParts[1])

	if err != nil {
		m.Utils.Logger.Fatalf("Invalid redis address: %v", err)
	}

	return redis.New(redis.Config{
		Host:     redisAddrParts[0],
		Port:     redisPort,
		Password: m.EnvVars.REDIS_PASSWORD,
		Database: m.EnvVars.REDIS_DB,
	})
}

func (m *Middlewares) VisitorLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        1,
		Storage:    m.getStorageConfig(),
		Expiration: time.Hour * 24,
		LimitReached: func(c fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "Too many requests")
		},
		LimiterMiddleware: limiter.FixedWindow{},
	})
}
