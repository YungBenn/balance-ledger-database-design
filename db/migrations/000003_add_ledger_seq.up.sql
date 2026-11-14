-- Per-user sequence number: the entry with the highest seq is the latest one.
-- created_at can't be used for this because now() is the transaction start time,
-- not the order in which concurrent writers actually commit.
ALTER TABLE ledger ADD COLUMN seq bigint;

UPDATE ledger l
SET seq = t.rn
FROM (
        SELECT id,
            row_number() OVER (PARTITION BY user_id ORDER BY created_at, id) AS rn
        FROM ledger
    ) t
WHERE l.id = t.id;

ALTER TABLE ledger ALTER COLUMN seq SET NOT NULL;

ALTER TABLE ledger ADD CONSTRAINT ledger_user_id_seq_key UNIQUE (user_id, seq);
