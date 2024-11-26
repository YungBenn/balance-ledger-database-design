-- name: ListLedgerByUser :many
SELECT *
FROM ledger
WHERE user_id = $1
LIMIT $3 OFFSET (($2 - 1) * $3);

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
        final
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
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
ORDER BY created_at DESC
LIMIT 1;
