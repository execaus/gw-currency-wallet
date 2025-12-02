package repository

import (
	"gw-currency-wallet/internal/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresRepository(pool *pgxpool.Pool) (*Repository, error) {
	queries := db.New(pool)

	return &Repository{
		Account: NewAccountRepository(pool, queries),
		Wallet:  NewWalletRepository(pool, queries),
	}, nil
}
