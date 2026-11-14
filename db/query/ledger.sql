-- name: ListLedgerByUser :many
SELECT *
FROM ledger
WHERE user_id = $1
ORDER BY seq DESC
LIMIT sqlc.arg(page_size) OFFSET sqlc.arg(page_offset);

-- name: GetLedgerByID :one
SELECT *
FROM ledger
WHERE id = $1
LIMIT 1;

-- name: CreateLedger :one
INSERT INTO ledger (
        id,
        user_id,
        TYPE,
        description,
        current,
        ADD,
        final,
        seq,
        idempotency_key
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetBalanceByUser :one
SELECT final
FROM ledger
WHERE user_id = $1
ORDER BY seq DESC
LIMIT 1;

-- name: GetLatestLedgerByUser :one
SELECT *
FROM ledger
WHERE user_id = $1
ORDER BY seq DESC
LIMIT 1;

-- name: GetLedgerByIdempotencyKey :one
SELECT *
FROM ledger
WHERE user_id = $1
    AND idempotency_key = $2
LIMIT 1;
