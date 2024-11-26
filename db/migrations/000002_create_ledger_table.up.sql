CREATE TYPE type_enum AS ENUM ('deposit', 'withdrawal', 'refund', 'purchase');

CREATE TABLE ledger (
    id varchar PRIMARY KEY,
    user_id varchar NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type type_enum NOT NULL,
    description varchar NOT NULL,
    current bigint NOT NULL,
    add bigint NOT NULL,
    final bigint NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now()),
    updated_at timestamptz NOT NULL DEFAULT (now())
);