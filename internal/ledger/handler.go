package ledger

import (
	"balance-ledger-database-design/db/sqlc"
	"balance-ledger-database-design/pkg/response"
	"balance-ledger-database-design/pkg/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type LedgerHandler struct {
	usecase LedgerUsecase
}

func NewLedgerHandler(usecase LedgerUsecase) *LedgerHandler {
	return &LedgerHandler{usecase}
}

func (h *LedgerHandler) CreateLedger(c *fiber.Ctx) error {
	var arg sqlc.CreateLedgerParams
	if err := c.BodyParser(&arg); err != nil {
		return err
	}

	res, err := h.usecase.CreateLedger(c.Context(), arg)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  http.StatusCreated,
		"message": "success",
		"ledger":  res,
	})
}

func (h *LedgerHandler) GetBalance(c *fiber.Ctx) error {
	user := utils.GetCurrentUser(c)

	res, err := h.usecase.GetBalance(c.Context(), user.UserID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "success",
		"balance": res,
	})
}

func (h *LedgerHandler) GetListLedger(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	limit := c.QueryInt("limit")
	offset := c.QueryInt("offset")

	res, err := h.usecase.GetListLedger(c.Context(), userID, int32(limit), int32(offset))
	if err != nil {
		return response.ErrorHandler(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "success",
		"ledgers": res,
	})
}

func (h *LedgerHandler) GetLedgerByID(c *fiber.Ctx) error {
	id := c.Params("id")

	res, err := h.usecase.GetLedgerByID(c.Context(), id)
	if err != nil {
		return response.ErrorHandler(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "success",
		"ledger":  res,
	})
}
