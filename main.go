package main

import (
	"context"
	"errors"

	"github.com/doz-8108/api-gateway/config"
	"github.com/doz-8108/api-gateway/db"
	"github.com/doz-8108/api-gateway/handlers"
	"github.com/doz-8108/api-gateway/middlewares"
	"github.com/doz-8108/api-gateway/storage"
	"github.com/doz-8108/api-gateway/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
	"go.uber.org/fx"
)

func newFiberServer(lc fx.Lifecycle, userHandlers *handlers.UserHandlers, visitorHandlers *handlers.VisitorHandlers, utils utils.Utils, middlewares middlewares.Middlewares, envVars config.EnvVars) *fiber.App {

	app := fiber.New(fiber.Config{ErrorHandler: func(f fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		message := "internal server error"

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) && fiberErr.Code != fiber.StatusInternalServerError {
			code = fiberErr.Code
			message = fiberErr.Message
		}

		if err != nil {
			return f.Status(code).JSON(fiber.Map{
				"message": message,
			})
		}

		return nil
	}})

	app.Use(recover.New())
	app.Get("/", static.New("./static/index.html"))

	visitorGroup := app.Group("/visitors")
	visitorGroup.Get("/", visitorHandlers.GetVisitorCounts)
	visitorGroup.Use(middlewares.VisitorLimiter()).Post("/", visitorHandlers.IncrementVisitorCount)

	// userGroup := app.Group("/users") // userGroup.Post("/sign-up", userHandlers.SignUpUser)
	// userGroup.Post("/sign-in", userHandlers.SignInUser)
	// userGroup.Get("/verify", userHandlers.VerifyUserEmail)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go app.Listen(":443", fiber.ListenConfig{
				CertFile:    envVars.TLS_CERT_FILE_PATH,
				CertKeyFile: envVars.TLS_KEY_FILE_PATH,
			})
			go app.Listen(":80")
			utils.Logger.Info("Server is listening at port 80 and port 443!")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			utils.Logger.Sync()
			return app.Shutdown()
		},
	})

	return app
}

func main() {
	fx.New(fx.Provide(config.LoadEnv, db.CreateMySqlConnection, db.CreateRedisConnection, utils.NewUtils,
		handlers.NewUserHandlers, handlers.NewVisitorHandlers, storage.NewUserStorage, middlewares.NewMiddlewares), fx.Invoke(newFiberServer)).Run()
}
