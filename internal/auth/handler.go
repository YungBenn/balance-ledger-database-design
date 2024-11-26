package auth

import (
	"balance-ledger-database-design/db/sqlc"
	"balance-ledger-database-design/pkg/response"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	usecase AuthUsecase
}

func NewAuthHandler(usecase AuthUsecase) *AuthHandler {
	return &AuthHandler{usecase}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var arg LoginRequest
	if err := c.BodyParser(&arg); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	res, err := h.usecase.Login(c.Context(), arg)
	if err != nil {
		return response.ErrorHandler(err)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "login success",
		"data":    res,
	})
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var arg sqlc.CreateUserParams
	if err := c.BodyParser(&arg); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	res, err := h.usecase.Register(c.Context(), arg)
	if err != nil {
		return response.ErrorHandler(err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  http.StatusCreated,
		"message": "register success",
		"data":    res,
	})
}
