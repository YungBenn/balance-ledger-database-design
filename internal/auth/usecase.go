package auth

import (
	"balance-ledger-database-design/db/sqlc"
	"balance-ledger-database-design/internal/postgresql"
	"balance-ledger-database-design/internal/token"
	"balance-ledger-database-design/pkg/response"
	"balance-ledger-database-design/pkg/utils"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type AuthUsecase interface {
	Register(ctx context.Context, arg sqlc.CreateUserParams) (*UserData, error)
	Login(ctx context.Context, arg LoginRequest) (*LoginResponse, error)
}

type authUsecase struct {
	repo postgresql.Repository
}

func NewAuthUsecase(repo postgresql.Repository) AuthUsecase {
	return &authUsecase{repo}
}

// Login implements AuthUsecase.
func (a *authUsecase) Login(ctx context.Context, arg LoginRequest) (*LoginResponse, error) {
	user, err := a.repo.GetUserByEmail(ctx, arg.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, response.NewError(http.StatusBadRequest, ErrInvalidCredentials)
		}
		return nil, response.NewError(http.StatusInternalServerError, err)
	}

	err = utils.CheckPassword(arg.Password, user.Password)
	if err != nil {
		return nil, response.NewError(http.StatusBadRequest, ErrInvalidCredentials)
	}

	token, _, err := token.Create(user.ID, user.Email, user.FullName, 1*time.Hour)
	if err != nil {
		return nil, response.NewError(http.StatusInternalServerError, err)
	}

	userData := UserData{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	loginResponse := &LoginResponse{
		Token: token,
		User:  userData,
	}

	return loginResponse, nil
}

// Register implements AuthUsecase.
func (a *authUsecase) Register(ctx context.Context, arg sqlc.CreateUserParams) (*UserData, error) {
	_, err := a.repo.GetUserByEmail(ctx, arg.Email)
	if err == nil {
		return nil, response.NewError(http.StatusBadRequest, ErrEmailExist)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, response.NewError(http.StatusInternalServerError, err)
	}

	hashedPass, errHashed := utils.HashPassword(arg.Password)
	if errHashed != nil {
		return nil, response.NewError(http.StatusInternalServerError, errHashed)
	}

	id, err := utils.GenerateCUID()
	if err != nil {
		return nil, response.NewError(http.StatusInternalServerError, err)
	}

	newUser := sqlc.CreateUserParams{
		ID:       id,
		Email:    arg.Email,
		FullName: arg.FullName,
		Password: hashedPass,
	}

	user, err := a.repo.CreateUser(ctx, newUser)
	if err != nil {
		// Two sign-ups with the same email can both pass the check above;
		// the UNIQUE constraint on email rejects the second one.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, response.NewError(http.StatusBadRequest, ErrEmailExist)
		}
		return nil, response.NewError(http.StatusInternalServerError, err)
	}

	// UserData has no password field, so the hash never reaches the client.
	return &UserData{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
