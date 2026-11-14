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
        seq
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateLedger :one
UPDATE ledger
SET user_id = $2,
    TYPE = $3,
    description = $4,
    current = $5,
    ADD = $6,
    final = $7,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteLedger :exec
DELETE FROM ledger
WHERE id = $1;

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
