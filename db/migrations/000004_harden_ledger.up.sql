-- Let Postgres enforce the ledger rules, not just the app.
-- If existing rows break a rule, this migration fails instead of hiding them.
ALTER TABLE ledger
    ADD CONSTRAINT ledger_final_check CHECK (final = current + add),
    ADD CONSTRAINT ledger_final_non_negative CHECK (final >= 0);

-- Optional client-chosen key: a retried request with the same key returns the
-- original entry instead of creating a second one.
ALTER TABLE ledger ADD COLUMN idempotency_key varchar;

ALTER TABLE ledger
    ADD CONSTRAINT ledger_user_id_idempotency_key_key UNIQUE (user_id, idempotency_key);

-- Entries are append-only: mistakes are corrected with a new entry, never edited.
-- DELETE stays allowed so deleting a user still cascades to their entries.
ALTER TABLE ledger DROP COLUMN updated_at;

CREATE FUNCTION ledger_reject_update() RETURNS trigger AS $$
BEGIN
    RAISE EXCEPTION 'ledger entries are append-only and cannot be updated';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER ledger_no_update
BEFORE UPDATE ON ledger
FOR EACH ROW EXECUTE FUNCTION ledger_reject_update();
