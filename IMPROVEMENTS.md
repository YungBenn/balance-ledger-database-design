# Improvements Checklist

Bugs, improvements, and refactors found while reviewing this learning project. Findings come from reading the code; each bug should be reproduced against a running instance before fixing it.

## Bugs

### Critical

- [x] **Password hash leaks on register.** `Register` returns the raw `sqlc.User`, whose `Password` field has a `json:"password"` tag (`db/sqlc/models.go:73`), so the bcrypt hash is sent back to the client.
- [x] **Lost update race in `CreateLedger`.** `internal/ledger/usecase.go` reads the latest balance and then inserts in two separate steps, with no transaction or lock. Two simultaneous requests for the same user both read the same `current`, and one change is lost.
- [x] **Down migration corrupts the catalog.** `db/migrations/000002_create_ledger_table.down.sql` deletes the enum from `pg_type` while `ledger` still uses it, then drops the table. Fix: `DROP TABLE ledger;` then `DROP TYPE type_enum;`.
- [x] **No authorization on ledger endpoints.** `POST /ledgers` trusts `user_id` from the body, and `GET /ledgers?user_id=` and `GET /ledgers/:id` need no login, so anyone can read or write anyone's ledger.

### Logic

- [x] **`/balance` is unreachable.** `/:id` is registered before `/balance` (`internal/ledger/routes.go:20-21`), so `GET /ledgers/balance` is handled as a ledger lookup with `id = "balance"`.
- [x] **Errors become an empty 200.** `response.ErrorHandler` returns `nil` for anything that isn't a `*response.Error` (`pkg/response/error.go:26`). For example, a database failure in `GetLedgerByID` (`internal/ledger/usecase.go:55`) makes the handler `return nil`.
- [x] **Pagination is broken.**
  - [x] The `offset` param is actually a page number: `OFFSET (page - 1) * limit`.
  - [x] `offset=0` produces a negative OFFSET, which Postgres rejects.
  - [x] With no query params, `limit` defaults to 0, so the list is always empty.
  - [x] `ListLedgerByUser` has no `ORDER BY`, so the order between pages isn't guaranteed.
- [x] **Balance can go negative.** Nothing rejects a withdrawal larger than the current balance.
- [x] **The sign of `add` isn't checked against `type`.** A `deposit` of -500 or a `withdrawal` of +500 is accepted.
- [x] **"Latest balance" can be ambiguous.** `GetBalanceByUser` orders by `created_at`, so two entries with the same timestamp give an unpredictable result.
- [x] **Every `CreateLedger` failure is a 500.** It returns plain `fmt.Errorf` errors, so an unknown user, which should be a 400 or 404, also comes back as a 500.
- [x] **A bad `type` value is a 500 instead of a 400.** It isn't validated before Postgres rejects it.
- [x] **Registration ignores database errors.** `Register` discards the error from `GetUserByEmail` (`internal/auth/usecase.go:68`). Two simultaneous sign-ups with the same email both pass the check, and the unique violation comes back as a 500.
- [x] **A failed ID generation is silently ignored.** `GenerateCUID` returns `""` on error (`pkg/utils/helper.go:15`), which can insert a row with an empty ID.
- [x] **`GetCurrentUser` can panic.** It assumes the user is always set, so it panics on any route without the auth middleware.
- [x] **Startup reports a connection it hasn't made.** `pgxpool.New` doesn't connect, so `postgresql.Connect` logs success even when the database is down. Add `db.Ping(ctx)`.
- [x] **The database pool is never closed.** `log.Fatal` exits immediately, so the deferred `postgresql.Close(db)` never runs.
- [x] **Several Makefile targets don't work.**
  - [x] `sqlc` uses a hardcoded `/home/rubenadi/...` path.
  - [x] `docker.stop` passes an image name, but `docker stop` needs a container name or ID.
  - [x] `docker.run` sets `DB_HOST=localhost`, which inside the container points at the container itself.

## Improvements: database design

- [x] Make ledger rows permanent: remove the `UpdateLedger` and `DeleteLedger` queries and fix mistakes with a new reversing entry. (Also a trigger rejects `UPDATE` in Postgres.)
- [x] Add `CHECK (final = current + add)`.
- [x] Add `CHECK (final >= 0)`.
- [x] Add a `CHECK` that the sign of `add` matches `type` (migration `000005`).
- [x] Add a per-user sequence number with `UNIQUE (user_id, seq)`. This settles which entry is latest and blocks the lost-update race.
- [x] Or lock the user row (`SELECT ... FROM users WHERE id = $1 FOR UPDATE`) inside a transaction while writing an entry. (Did both: the lock serializes writers, the unique `seq` is the safety net.)
- [x] Add an index on `ledger (user_id, created_at DESC)`, or on `(user_id, seq)` if you add the sequence. (Covered by the `UNIQUE (user_id, seq)` index.)
- [x] Add an idempotency key so a retried request doesn't create the same entry twice.
- [x] Drop `updated_at` from `ledger` once rows are no longer edited.

## Refactors: code

- [x] **Keep one sqlc output folder.** `sqlc.yaml` writes to `internal/postgresql/sqlc`, but the app imports `db/sqlc`, so `make sqlc` updates code nobody uses.
- [x] **Name the pagination params** with `sqlc.arg(page)` and `sqlc.arg(page_size)` instead of the generated `Column2 interface{}`.
- [x] **Add transaction support to the repository.** `SQLRepository` already holds the `db` pool but never uses it.
- [x] **Use request and response types instead of `sqlc` types in handlers.** `internal/ledger/dto.go` is empty. This also fixes the password-hash leak.
- [x] **Validate input at the handler.**
  - [x] Email format
  - [x] Password length
  - [x] `add != 0`
  - [x] `type` is one of the allowed values
- [x] **Use one error path.** Every usecase returns `response.Error` with the right status, and the global Fiber error handler maps it.
- [x] **Make response shapes consistent.** The top-level key is currently one of `data`, `ledger`, `ledgers` or `balance`.
- [x] **Remove the dead line in `FiberErrorHandler`** that sets text/plain right before writing JSON.
- [x] **Load the PASETO key from the environment** instead of hardcoding it in `internal/token/paseto.go`.
- [x] **Add graceful shutdown** with `app.Shutdown` on SIGINT/SIGTERM.
- [x] **Pin the Go version in the Dockerfile** (`golang:1.22-alpine`).
- [x] **Add tests.**
  - [x] Integration test for `CreateLedger` that checks the balance chain stays consistent (`row[n].current = row[n-1].final`).
  - [x] Concurrency test that sends simultaneous `CreateLedger` requests to catch the lost-update race.
