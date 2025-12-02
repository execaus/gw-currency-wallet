package repository

import (
	"context"
	"errors"
	"gw-currency-wallet/internal/db"
	"gw-currency-wallet/pkg"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type WalletRepository struct {
	TxRepositoryImpl
}

func (r *WalletRepository) IsExistCurrency(ctx context.Context, email string, currency pkg.Currency) (bool, error) {
	q := r.getQueries(ctx)

	result, err := q.IsExistCurrency(ctx, db.IsExistCurrencyParams{
		Email:    email,
		Currency: currency,
	})
	if err != nil {
		zap.L().Error(err.Error())
		return false, err
	}

	return result, err
}

func (r *WalletRepository) Create(ctx context.Context, email string, currency pkg.Currency) error {
	q := r.getQueries(ctx)

	if err := q.CreateWallet(ctx, db.CreateWalletParams{
		Email:    email,
		Currency: currency,
	}); err != nil {
		zap.L().Error(err.Error())
		return err
	}

	return nil
}

func (r *WalletRepository) GetForUpdate(ctx context.Context, email string, currency pkg.Currency) (*db.AppWallet, error) {
	q := r.getQueries(ctx)

	row, err := q.GetWalletForUpdate(ctx, db.GetWalletForUpdateParams{
		Email:    email,
		Currency: currency,
	})
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, nil
		default:
			zap.L().Error(err.Error())
			return nil, err
		}
	}

	return &row, nil
}

func (r *WalletRepository) Update(ctx context.Context, email string, currency pkg.Currency, newValue float32) (*db.AppWallet, error) {
	q := r.getQueries(ctx)

	row, err := q.UpdateWallet(ctx, db.UpdateWalletParams{
		Email:    email,
		Currency: currency,
		Balance:  newValue,
	})
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	return &row, nil
}

func (r *WalletRepository) GetAllByEmail(ctx context.Context, email string) ([]db.AppWallet, error) {
	q := r.getQueries(ctx)

	rows, err := q.GetWalletsByEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, nil
		default:
			zap.L().Error(err.Error())
			return nil, err
		}
	}

	return rows, nil
}

func NewWalletRepository(pool *pgxpool.Pool, queries *db.Queries) *WalletRepository {
	return &WalletRepository{
		TxRepositoryImpl{
			db: pool,
			q:  queries,
		},
	}
}
