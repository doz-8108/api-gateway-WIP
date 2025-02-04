package middlewares

type Middlewares struct {
	JWT_SIGN_SECRET string
}

func NewMiddlewares(jwtSignSecret string) Middlewares {
	return Middlewares{
		JWT_SIGN_SECRET: jwtSignSecret,
	}
}
