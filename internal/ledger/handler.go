package ledger

import (
	"balance-ledger-database-design/db/sqlc"
	"balance-ledger-database-design/pkg/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
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

	var req CreateLedgerRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	key := c.Get("Idempotency-Key")
	arg := sqlc.CreateLedgerParams{
		UserID:         user.UserID, // always the logged-in user
		Type:           sqlc.TypeEnum(req.Type),
		Description:    req.Description,
		Add:            req.Add,
		IdempotencyKey: pgtype.Text{String: key, Valid: key != ""},
	}

	res, err := h.usecase.CreateLedger(c.UserContext(), arg)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  http.StatusCreated,
		"message": "success",
		"data":    toLedgerResponse(*res),
	})
}

func (h *LedgerHandler) GetBalance(c *fiber.Ctx) error {
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		return err
	}

	res, err := h.usecase.GetBalance(c.UserContext(), user.UserID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "success",
		"data":    BalanceResponse{Balance: *res},
	})
}

func (h *LedgerHandler) GetListLedger(c *fiber.Ctx) error {
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		return err
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	res, err := h.usecase.GetListLedger(c.UserContext(), user.UserID, int32(page), int32(limit))
	if err != nil {
		return err
	}

	ledgers := make([]LedgerResponse, len(res))
	for i, l := range res {
		ledgers[i] = toLedgerResponse(l)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "success",
		"data":    ledgers,
	})
}

func (h *LedgerHandler) GetLedgerByID(c *fiber.Ctx) error {
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		return err
	}

	res, err := h.usecase.GetLedgerByID(c.UserContext(), c.Params("id"), user.UserID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "success",
		"data":    toLedgerResponse(*res),
	})
}
