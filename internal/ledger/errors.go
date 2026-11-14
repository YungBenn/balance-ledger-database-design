package ledger

import "errors"

var (
	ErrLedgerNotFound      = errors.New("ledger not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidType         = errors.New("type must be one of deposit, withdrawal, refund, purchase")
	ErrZeroAmount          = errors.New("add must not be zero")
	ErrInvalidSign         = errors.New("add must be positive for deposit and refund, negative for withdrawal and purchase")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidPagination   = errors.New("page must be at least 1 and limit between 1 and 100")
	ErrInvalidIdemKey      = errors.New("Idempotency-Key must be at most 255 characters")
	ErrIdemKeyReused       = errors.New("Idempotency-Key was already used for a different request")
)
