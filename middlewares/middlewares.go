package middlewares

import (
	"github.com/doz-8108/api-gateway/config"
	"github.com/doz-8108/api-gateway/utils"
)

type Middlewares struct {
	EnvVars config.EnvVars
	Utils   utils.Utils
}

func NewMiddlewares(envVars config.EnvVars, utils utils.Utils) Middlewares {
	return Middlewares{
		EnvVars: envVars,
		Utils:   utils,
	}
}
