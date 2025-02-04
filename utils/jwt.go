package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (u *Utils) CreateToken(userId int, issuer string) (string, error) {
	currentTime := time.Now()
	encodedUserId := u.EncodeUserId(uint64(userId))
	claims :=
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Hour * 24 * 7)),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			NotBefore: jwt.NewNumericDate(currentTime),
			Issuer:    issuer,
			Subject:   encodedUserId,
		}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(u.EnvVars.JWT_SIGN_SECRET))
}
