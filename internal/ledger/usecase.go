package ledger

import (
	"balance-ledger-database-design/db/sqlc"
	"balance-ledger-database-design/internal/postgresql"
	"balance-ledger-database-design/pkg/response"
	"balance-ledger-database-design/pkg/utils"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
)

// LedgerUsecase defines the interface for ledger business logic operations.
// It provides methods to create ledger entries and retrieve balance information.
type LedgerUsecase interface {
	// CreateLedger creates a new ledger entry for a user with the given parameters.
	// Returns the created ledger entry and any error encountered during the process.
	CreateLedger(ctx context.Context, arg sqlc.CreateLedgerParams) (*sqlc.Ledger, error)

	// GetBalance retrieves the current balance for a specific user.
	// Returns the balance amount and any error encountered during the process.
	GetBalance(ctx context.Context, userID string) (*int64, error)

	// GetListLedger retrieves a list of ledger entries for a specific user.
	// Returns the list of ledger entries and any error encountered during the process.
	GetListLedger(ctx context.Context, userID string, limit int32, offset int32) ([]sqlc.Ledger, error)

	// GetLedgerByID retrieves a ledger entry by its ID.
	// Returns the ledger entry and any error encountered during the process.
	GetLedgerByID(ctx context.Context, ledgerID string) (*sqlc.Ledger, error)
}

// ledgerUsecase implements the LedgerUsecase interface.
// It uses a postgresql.Repository for data persistence.
type ledgerUsecase struct {
	repo postgresql.Repository
}

// NewLedgerUsecase creates a new instance of LedgerUsecase.
// It takes a postgresql.Repository as a dependency for data operations.
func NewLedgerUsecase(repo postgresql.Repository) LedgerUsecase {
	return &ledgerUsecase{repo}
}

// GetLedgerByID implements LedgerUsecase.
func (a *ledgerUsecase) GetLedgerByID(ctx context.Context, ledgerID string) (*sqlc.Ledger, error) {
	ledger, err := a.repo.GetLedgerByID(ctx, ledgerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, response.NewError(http.StatusNotFound, ErrLedgerNotFound)
		}
		return nil, err
	}

	return &ledger, nil
}

// GetListLedger implements LedgerUsecase.
func (a *ledgerUsecase) GetListLedger(ctx context.Context, userID string, limit int32, offset int32) ([]sqlc.Ledger, error) {
	arg := sqlc.ListLedgerByUserParams{
		UserID:  userID,
		Column2: offset,
		Limit:   limit,
	}

	ledgers, err := a.repo.ListLedgerByUser(ctx, arg)
	if err != nil {
		return nil, response.NewError(http.StatusInternalServerError, err)
	}

	return ledgers, nil
}

// GetBalance retrieves the current balance for a specified user.
// It queries the repository to get the user's balance and returns it.
// If any error occurs during the process, it returns nil and the error.
func (a *ledgerUsecase) GetBalance(ctx context.Context, userID string) (*int64, error) {
	balance, err := a.repo.GetBalanceByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zero := int64(0)
			return &zero, nil
		}
		return nil, err
	}

	return &balance, nil
}

// validateUser checks if the user exists in the system
func (a *ledgerUsecase) validateUser(ctx context.Context, userID string) error {
	if _, err := a.repo.GetUserByID(ctx, userID); err != nil {
		return fmt.Errorf("invalid user: %w", err)
	}
	return nil
}

func (a *ledgerUsecase) CreateLedger(ctx context.Context, arg sqlc.CreateLedgerParams) (*sqlc.Ledger, error) {
	// Validate user existence
	if err := a.validateUser(ctx, arg.UserID); err != nil {
		return nil, err
	}

	// Get current balance, defaulting to 0 if not found
	currentBalance, err := a.repo.GetBalanceByUser(ctx, arg.UserID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	// Prepare new ledger entry
	newLedger := sqlc.CreateLedgerParams{
		ID:          utils.GenerateCUID(),
		UserID:      arg.UserID,
		Type:        arg.Type,
		Description: arg.Description,
		Current:     currentBalance,
		Add:         arg.Add,
		Final:       currentBalance + arg.Add,
	}

	// Create ledger entry in database
	ledgerData, err := a.repo.CreateLedger(ctx, newLedger)
	if err != nil {
		return nil, fmt.Errorf("failed to create ledger: %w", err)
	}

	return &ledgerData, nil
}
