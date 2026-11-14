package main

import (
	"balance-ledger-database-design/internal/auth"
	"balance-ledger-database-design/internal/ledger"
	"balance-ledger-database-design/internal/postgresql"
	"balance-ledger-database-design/pkg/utils"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)

	db, err := postgresql.Connect(context.Background(), dbURL)
	if err != nil {
		log.Panic(err)
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          utils.FiberErrorHandler,
	})

	auth.SetupAuthRoutes(app, db)
	ledger.SetupLedgerRoutes(app, db)

	log.Printf("Server started on http://127.0.0.1:%s", os.Getenv("PORT"))
	err = app.Listen(":" + os.Getenv("PORT"))
	// Close explicitly: log.Fatal exits without running deferred calls.
	postgresql.Close(db)
	if err != nil {
		log.Fatal(err)
	}
}
