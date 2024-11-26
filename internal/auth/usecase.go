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
)

type AuthUsecase interface {
	Register(ctx context.Context, arg sqlc.CreateUserParams) (*sqlc.User, error)
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
func (a *authUsecase) Register(ctx context.Context, arg sqlc.CreateUserParams) (*sqlc.User, error) {
	user, _ := a.repo.GetUserByEmail(ctx, arg.Email)
	if user.Email != "" {
		return nil, response.NewError(http.StatusBadRequest, ErrEmailExist)
	}

	hashedPass, errHashed := utils.HashPassword(arg.Password)
	if errHashed != nil {
		return nil, response.NewError(http.StatusInternalServerError, errHashed)
	}

	newUser := sqlc.CreateUserParams{
		ID:       utils.GenerateCUID(),
		Email:    arg.Email,
		FullName: arg.FullName,
		Password: hashedPass,
	}

	user, err := a.repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, response.NewError(http.StatusInternalServerError, err)
	}

	return &user, nil
}
