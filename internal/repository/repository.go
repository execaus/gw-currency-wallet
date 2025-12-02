package repository

import (
	"context"
	"gw-currency-wallet/internal/db"
	"gw-currency-wallet/pkg"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Wallet interface {
	TxRepository
	GetAllByEmail(ctx context.Context, email string) ([]db.AppWallet, error)
	GetForUpdate(ctx context.Context, email string, currency pkg.Currency) (*db.AppWallet, error)
	Update(ctx context.Context, email string, currency pkg.Currency, newValue float32) (*db.AppWallet, error)
	IsExistCurrency(ctx context.Context, email string, currency pkg.Currency) (bool, error)
	Create(ctx context.Context, email string, currency pkg.Currency) error
}

type Account interface {
	TxRepository
	IsEmailExist(ctx context.Context, email string) (bool, error)
	IsUsernameExist(ctx context.Context, username string) (bool, error)
	Create(ctx context.Context, email, username, passwordHash string) (*db.AppAccount, error)
	GetByUsername(ctx context.Context, username string) (*db.AppAccount, error)
}

type Repository struct {
	Wallet
	Account
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	r, err := NewPostgresRepository(pool)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	return r
}
