package utils

import "github.com/gofiber/fiber/v3"

func (u *Utils) CatchError(err any, httpErr *fiber.Error) {
	if err != nil {
		u.Logger.Error(err)
	}

	if httpErr != nil {
		panic(httpErr)
	}
}
