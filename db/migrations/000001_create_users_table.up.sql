CREATE TABLE users (
    id varchar PRIMARY KEY,
    email varchar UNIQUE NOT NULL,
    full_name varchar NOT NULL,
    password varchar NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now()),
    updated_at timestamptz NOT NULL DEFAULT (now())
);
