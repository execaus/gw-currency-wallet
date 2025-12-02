package repository

import (
	"context"
	"errors"
	"gw-currency-wallet/internal/db"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type WalletRepository struct {
	q *db.Queries
}

func (r *WalletRepository) GetAllByEmail(ctx context.Context, email string) ([]db.AppWallet, error) {
	rows, err := r.q.GetWalletsByEmail(ctx, email)
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

func NewWalletRepository(q *db.Queries) *WalletRepository {
	return &WalletRepository{
		q: q,
	}
}
