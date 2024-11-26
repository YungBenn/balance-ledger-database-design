package postgresql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	log.Println("Success connect to database")

	return db, nil
}

func Close(db *pgxpool.Pool) {
	db.Close()
	log.Println("Success close connection to database")
}
