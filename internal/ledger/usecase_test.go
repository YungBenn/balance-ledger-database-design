package ledger

import (
	"balance-ledger-database-design/db/sqlc"
	"balance-ledger-database-design/internal/postgresql"
	"balance-ledger-database-design/pkg/utils"
	"context"
	"os"
	"sync"
	"testing"
)

// Concurrent entries for one user must chain without losing updates.
// Needs a migrated database: TEST_DATABASE_URL=postgres://... go test ./...
func TestCreateLedgerConcurrent(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}

	ctx := context.Background()
	db, err := postgresql.Connect(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := postgresql.NewRepository(db)
	user, err := repo.CreateUser(ctx, sqlc.CreateUserParams{
		ID:       utils.GenerateCUID(),
		Email:    utils.GenerateCUID() + "@test.local",
		FullName: "Race Test",
		Password: "x",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer repo.DeleteUser(ctx, user.ID)

	const n = 50
	uc := NewLedgerUsecase(repo)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := uc.CreateLedger(ctx, sqlc.CreateLedgerParams{
				UserID: user.ID, Type: sqlc.TypeEnumDeposit, Description: "race", Add: 1,
			})
			if err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()

	balance, err := uc.GetBalance(ctx, user.ID)
	if err != nil {
		t.Fatal(err)
	}
	if *balance != n {
		t.Fatalf("balance = %d, want %d (lost updates)", *balance, n)
	}
}
