package auth

import (
	"balance-ledger-database-design/internal/postgresql"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupAuthRoutes(app *fiber.App, db *pgxpool.Pool) {
	repo := postgresql.NewRepository(db)
	usecase := NewAuthUsecase(repo)
	handler := NewAuthHandler(usecase)

	r := app.Group("/api/v1/auth")

	r.Post("/register", handler.Register)
	r.Post("/login", handler.Login)
}
