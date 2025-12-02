package repository

import (
	"context"
	"gw-currency-wallet/internal/db"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type AccountRepository struct {
	TxRepositoryImpl
}

func (r *AccountRepository) GetByUsername(ctx context.Context, username string) (*db.AppAccount, error) {
	q := r.getQueries(ctx)

	account, err := q.GetAccountByUsername(ctx, username)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) IsEmailExist(ctx context.Context, email string) (bool, error) {
	q := r.getQueries(ctx)

	isExist, err := q.IsAccountExistsByEmail(ctx, email)
	if err != nil {
		zap.L().Error(err.Error())
		return false, err
	}

	return isExist, nil
}

func (r *AccountRepository) IsUsernameExist(ctx context.Context, username string) (bool, error) {
	q := r.getQueries(ctx)

	isExist, err := q.IsAccountExistsByUsername(ctx, username)
	if err != nil {
		zap.L().Error(err.Error())
		return false, err
	}

	return isExist, nil
}

func (r *AccountRepository) Create(ctx context.Context, email, username, passwordHash string) (*db.AppAccount, error) {
	q := r.getQueries(ctx)

	row, err := q.CreateAccount(ctx, db.CreateAccountParams{
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

func NewAccountRepository(pool *pgxpool.Pool, queries *db.Queries) *AccountRepository {
	return &AccountRepository{
		TxRepositoryImpl{
			db: pool,
			q:  queries,
		},
	}
}
