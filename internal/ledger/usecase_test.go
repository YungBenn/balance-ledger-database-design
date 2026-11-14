package ledger

import (
	"balance-ledger-database-design/db/sqlc"
	"balance-ledger-database-design/internal/postgresql"
	"balance-ledger-database-design/pkg/response"
	"balance-ledger-database-design/pkg/utils"
	"context"
	"errors"
	"net/http"
	"os"
	"sync"
	"testing"
)

func TestValidateEntry(t *testing.T) {
	cases := []struct {
		typ  sqlc.TypeEnum
		add  int64
		want error
	}{
		{sqlc.TypeEnumDeposit, 100, nil},
		{sqlc.TypeEnumRefund, 5, nil},
		{sqlc.TypeEnumWithdrawal, -100, nil},
		{sqlc.TypeEnumPurchase, -5, nil},
		{sqlc.TypeEnumDeposit, -100, ErrInvalidSign},
		{sqlc.TypeEnumRefund, -5, ErrInvalidSign},
		{sqlc.TypeEnumWithdrawal, 100, ErrInvalidSign},
		{sqlc.TypeEnumPurchase, 5, ErrInvalidSign},
		{sqlc.TypeEnumDeposit, 0, ErrZeroAmount},
		{"bogus", 5, ErrInvalidType},
	}
	for _, c := range cases {
		err := validateEntry(c.typ, c.add)
		if !errors.Is(errUnwrap(err), c.want) {
			t.Errorf("validateEntry(%s, %d) = %v, want %v", c.typ, c.add, err, c.want)
		}
	}
}

// errUnwrap returns the inner error of a *response.Error so it can be compared with errors.Is.
func errUnwrap(err error) error {
	var re *response.Error
	if errors.As(err, &re) {
		return re.Err
	}
	return err
}

// newTestUser creates a user in TEST_DATABASE_URL and returns a usecase for it.
// Needs a migrated database: TEST_DATABASE_URL=postgres://... go test ./...
func newTestUser(t *testing.T) (LedgerUsecase, string) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}

	ctx := context.Background()
	db, err := postgresql.Connect(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(db.Close)

	repo := postgresql.NewRepository(db)
	id, err := utils.GenerateCUID()
	if err != nil {
		t.Fatal(err)
	}
	user, err := repo.CreateUser(ctx, sqlc.CreateUserParams{
		ID:       id,
		Email:    id + "@test.local",
		FullName: "Test",
		Password: "x",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { repo.DeleteUser(ctx, user.ID) })

	return NewLedgerUsecase(repo), user.ID
}

// Concurrent entries for one user must chain without losing updates.
func TestCreateLedgerConcurrent(t *testing.T) {
	uc, userID := newTestUser(t)
	ctx := context.Background()

	const n = 50
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := uc.CreateLedger(ctx, sqlc.CreateLedgerParams{
				UserID: userID, Type: sqlc.TypeEnumDeposit, Description: "race", Add: 1,
			})
			if err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()

	balance, err := uc.GetBalance(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}
	if *balance != n {
		t.Fatalf("balance = %d, want %d (lost updates)", *balance, n)
	}
}

func TestCreateLedgerRejectsOverdraft(t *testing.T) {
	uc, userID := newTestUser(t)
	ctx := context.Background()

	if _, err := uc.CreateLedger(ctx, sqlc.CreateLedgerParams{
		UserID: userID, Type: sqlc.TypeEnumDeposit, Description: "top up", Add: 100,
	}); err != nil {
		t.Fatal(err)
	}

	_, err := uc.CreateLedger(ctx, sqlc.CreateLedgerParams{
		UserID: userID, Type: sqlc.TypeEnumWithdrawal, Description: "too much", Add: -101,
	})
	var re *response.Error
	if !errors.As(err, &re) || re.Status != http.StatusBadRequest || !errors.Is(re.Err, ErrInsufficientBalance) {
		t.Fatalf("err = %v, want 400 %v", err, ErrInsufficientBalance)
	}

	balance, err := uc.GetBalance(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}
	if *balance != 100 {
		t.Fatalf("balance = %d, want 100", *balance)
	}
}
