package main

import (
	"balance-ledger-database-design/internal/auth"
	"balance-ledger-database-design/internal/ledger"
	"balance-ledger-database-design/internal/postgresql"
	"balance-ledger-database-design/internal/token"
	"balance-ledger-database-design/pkg/utils"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	if err := token.SetSymmetricKey(os.Getenv("TOKEN_SYMMETRIC_KEY")); err != nil {
		log.Fatal(err)
	}

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	listenErr := make(chan error, 1)
	go func() {
		listenErr <- app.Listen(":" + os.Getenv("PORT"))
	}()
	log.Printf("Server started on http://127.0.0.1:%s", os.Getenv("PORT"))

	// Run until the listener fails or SIGINT/SIGTERM arrives; then let in-flight requests finish.
	select {
	case err = <-listenErr:
	case <-ctx.Done():
		log.Println("Shutting down, waiting for in-flight requests")
		err = app.ShutdownWithTimeout(10 * time.Second)
	}

	// Close explicitly: log.Fatal exits without running deferred calls.
	postgresql.Close(db)
	if err != nil {
		log.Fatal(err)
	}
}
