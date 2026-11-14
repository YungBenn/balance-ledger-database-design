package postgresql

import (
	"balance-ledger-database-design/db/sqlc"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	sqlc.Querier
	// ExecTx runs fn inside a transaction, committing if fn returns nil and rolling back otherwise.
	ExecTx(ctx context.Context, fn func(q sqlc.Querier) error) error
}

type SQLRepository struct {
	db *pgxpool.Pool
	*sqlc.Queries
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &SQLRepository{
		db:      db,
		Queries: sqlc.New(db),
	}
}

func (r *SQLRepository) ExecTx(ctx context.Context, fn func(q sqlc.Querier) error) error {
	return pgx.BeginFunc(ctx, r.db, func(tx pgx.Tx) error {
		return fn(r.Queries.WithTx(tx))
	})
}
