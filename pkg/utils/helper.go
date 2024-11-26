package utils

import (
	"balance-ledger-database-design/internal/token"
	"crypto/rand"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/lucsky/cuid"
)

func GenerateCUID() string {
	cuid, err := cuid.NewCrypto(rand.Reader)
	if err != nil {
		return ""
	}

	return cuid
}

func GetCurrentUser(c *fiber.Ctx) *token.Payload {
	return c.Locals("user").(*token.Payload)
}

func FiberErrorHandler(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	// Set Content-Type: text/plain; charset=utf-8
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	// Return status code with error message
	return c.Status(code).JSON(fiber.Map{
		"status":  code,
		"message": err.Error(),
	})
}
