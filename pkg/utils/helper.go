package utils

import (
	"balance-ledger-database-design/internal/token"
	"balance-ledger-database-design/pkg/response"
	"crypto/rand"
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/lucsky/cuid"
)

func GenerateCUID() (string, error) {
	id, err := cuid.NewCrypto(rand.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to generate id: %w", err)
	}

	return id, nil
}

// GetCurrentUser returns the user set by the Authenticate middleware,
// or a 401 error if the route isn't behind it.
func GetCurrentUser(c *fiber.Ctx) (*token.Payload, error) {
	user, ok := c.Locals("user").(*token.Payload)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	return user, nil
}

// FiberErrorHandler is the single place errors become HTTP responses.
// It understands both *response.Error (from usecases) and *fiber.Error (from handlers and middleware).
func FiberErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := err.Error()

	var re *response.Error
	var fe *fiber.Error
	switch {
	case errors.As(err, &re):
		code, message = re.Status, re.Err.Error()
	case errors.As(err, &fe):
		code, message = fe.Code, fe.Message
	}

	// Internal errors can contain SQL or connection details, so log them instead of returning them.
	if code >= fiber.StatusInternalServerError {
		log.Printf("internal error on %s %s: %v", c.Method(), c.Path(), err)
		message = "internal server error"
	}

	return c.Status(code).JSON(fiber.Map{
		"status":  code,
		"message": message,
	})
}
