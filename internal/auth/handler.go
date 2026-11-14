package auth

import (
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

	res, err := h.usecase.Login(c.UserContext(), arg)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "login success",
		"data":    res,
	})
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var arg RegisterRequest
	if err := c.BodyParser(&arg); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	res, err := h.usecase.Register(c.UserContext(), arg)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  http.StatusCreated,
		"message": "register success",
		"data":    res,
	})
}
