package ledger

import (
	"balance-ledger-database-design/db/sqlc"
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
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		return err
	}

	var arg sqlc.CreateLedgerParams
	if err := c.BodyParser(&arg); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	// Always write to the logged-in user, never to a user_id from the body.
	arg.UserID = user.UserID

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
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		return err
	}

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
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		return err
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	res, err := h.usecase.GetListLedger(c.Context(), user.UserID, int32(page), int32(limit))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "success",
		"ledgers": res,
	})
}

func (h *LedgerHandler) GetLedgerByID(c *fiber.Ctx) error {
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		return err
	}

	res, err := h.usecase.GetLedgerByID(c.Context(), c.Params("id"), user.UserID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "success",
		"ledger":  res,
	})
}
