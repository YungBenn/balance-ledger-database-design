DROP TRIGGER IF EXISTS ledger_no_update ON ledger;

DROP FUNCTION IF EXISTS ledger_reject_update();

ALTER TABLE ledger ADD COLUMN updated_at timestamptz NOT NULL DEFAULT (now());

ALTER TABLE ledger
    DROP CONSTRAINT IF EXISTS ledger_user_id_idempotency_key_key,
    DROP COLUMN IF EXISTS idempotency_key,
    DROP CONSTRAINT IF EXISTS ledger_final_non_negative,
    DROP CONSTRAINT IF EXISTS ledger_final_check;
