package repository

import (
	"context"
	"gw-currency-wallet/internal/db"

	"go.uber.org/zap"
)

type AccountRepository struct {
	q *db.Queries
}

func (r *AccountRepository) GetByUsername(ctx context.Context, username string) (*db.AppAccount, error) {
	account, err := r.q.GetAccountByUsername(ctx, username)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) IsEmailExist(ctx context.Context, email string) (bool, error) {
	isExist, err := r.q.IsAccountExistsByEmail(ctx, email)
	if err != nil {
		zap.L().Error(err.Error())
		return false, err
	}

	return isExist, nil
}

func (r *AccountRepository) IsUsernameExist(ctx context.Context, username string) (bool, error) {
	isExist, err := r.q.IsAccountExistsByUsername(ctx, username)
	if err != nil {
		zap.L().Error(err.Error())
		return false, err
	}

	return isExist, nil
}

func (r *AccountRepository) Create(ctx context.Context, email, username, passwordHash string) (*db.AppAccount, error) {
	row, err := r.q.CreateAccount(ctx, db.CreateAccountParams{
		Email:    email,
		Username: username,
		Password: passwordHash,
	})
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	return &row, nil
}

func NewAccountRepository(queries *db.Queries) *AccountRepository {
	return &AccountRepository{
		q: queries,
	}
}
