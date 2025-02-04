package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/doz-8108/api-gateway/config"
	"github.com/doz-8108/api-gateway/db"
	"github.com/doz-8108/api-gateway/handlers"
	"github.com/doz-8108/api-gateway/storage"
	"github.com/doz-8108/api-gateway/utils"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
)

func newFiberServer(lc fx.Lifecycle, userHandlers *handlers.UserHandlers) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: func(f fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		message := "Internal Server Error"

		var e *fiber.Error
		if errors.As(err, &e) && e.Code != fiber.StatusInternalServerError {
			code = e.Code
			message = e.Message
		}

		// TODO: log error

		if err != nil {
			return f.Status(code).JSON(fiber.Map{
				"message": message,
			})
		}

		return nil
	}})

	// TODO: add middlewares

	userGroup := app.Group("/users")
	userGroup.Post("/sign-up", userHandlers.SignUpUser)
	userGroup.Post("/sign-in", userHandlers.SignInUser)
	userGroup.Get("/verify", userHandlers.VerifyUserEmail)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go app.Listen(":8080")
			fmt.Println("Server is listening at port 8080")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})

	return app
}

func main() {
	fx.New(fx.Provide(config.LoadEnv, db.CreateMySqlConnection, db.CreateRedisConnection, utils.NewUtils,
		handlers.NewUserHandlers, storage.NewUserStorage), fx.Invoke(newFiberServer)).Run()
}
