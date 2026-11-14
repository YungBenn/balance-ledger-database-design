package ledger

import (
	"balance-ledger-database-design/db/sqlc"
	"time"
)

// CreateLedgerRequest is the only input a client controls; current, final, seq,
// and user_id are always set by the server.
type CreateLedgerRequest struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Add         int64  `json:"add"`
}

type LedgerResponse struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Type           string    `json:"type"`
	Description    string    `json:"description"`
	Current        int64     `json:"current"`
	Add            int64     `json:"add"`
	Final          int64     `json:"final"`
	Seq            int64     `json:"seq"`
	IdempotencyKey *string   `json:"idempotency_key,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type BalanceResponse struct {
	Balance int64 `json:"balance"`
}

func toLedgerResponse(l sqlc.Ledger) LedgerResponse {
	res := LedgerResponse{
		ID:          l.ID,
		UserID:      l.UserID,
		Type:        string(l.Type),
		Description: l.Description,
		Current:     l.Current,
		Add:         l.Add,
		Final:       l.Final,
		Seq:         l.Seq,
		CreatedAt:   l.CreatedAt,
	}
	if l.IdempotencyKey.Valid {
		res.IdempotencyKey = &l.IdempotencyKey.String
	}
	return res
}
