package middlewares

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber"
	"github.com/golang-jwt/jwt/v5"
)

func (m *Middlewares) Auth(f fiber.Ctx) error {
	authorization := f.Get("Authorization")
	unAuthErr := fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")

	if !strings.HasPrefix(authorization, "Bearer ") {
		return unAuthErr
	}

	segs := strings.Split(authorization, "Bearer ")
	if len(segs) < 2 {
		return unAuthErr
	}

	token, err := jwt.Parse(segs[1], func(token *jwt.Token) (interface{}, error) {
		return m.EnvVars.JWT_SIGN_SECRET, nil
	})

	switch {
	case token.Valid:
		sub, err := token.Claims.GetSubject()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		f.Set("userId", sub)
		f.Next()
	case errors.Is(err, jwt.ErrTokenMalformed):
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
	default:
		return unAuthErr
	}

	return nil
}
