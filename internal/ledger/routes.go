package ledger

import (
	"balance-ledger-database-design/internal/middlewares"
	"balance-ledger-database-design/internal/postgresql"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupLedgerRoutes(app *fiber.App, db *pgxpool.Pool) {
	repo := postgresql.NewRepository(db)
	usecase := NewLedgerUsecase(repo)
	handler := NewLedgerHandler(usecase)

	r := app.Group("/api/v1/ledgers", middlewares.Authenticate())

	r.Post("/", handler.CreateLedger)
	r.Get("/", handler.GetListLedger)
	// Must be registered before /:id, or "balance" is matched as a ledger id.
	r.Get("/balance", handler.GetBalance)
	r.Get("/:id", handler.GetLedgerByID)
}
