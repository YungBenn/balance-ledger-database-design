package middlewares

import (
	"balance-ledger-database-design/internal/token"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(http.StatusUnauthorized, "missing authorization header")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		payload, err := token.Verify(tokenString)
		if err != nil {
			return fiber.NewError(http.StatusUnauthorized, err.Error())
		}

		c.Locals("user", payload)

		return c.Next()
	}
}
