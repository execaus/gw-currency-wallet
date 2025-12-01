package repository

import (
	"context"
	"gw-currency-wallet/internal/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Account interface {
	IsEmailExist(ctx context.Context, email string) (bool, error)
	IsUsernameExist(ctx context.Context, username string) (bool, error)
	Create(ctx context.Context, email, username, passwordHash string) (*db.AppAccount, error)
	GetByUsername(ctx context.Context, username string) (*db.AppAccount, error)
}

type Repository struct {
	Account
}

func NewRepository(pool *pgxpool.Pool) Repository {
	q := db.New(pool)

	return Repository{
		Account: NewAccountRepository(q),
	}
}
