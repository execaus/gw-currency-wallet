package service

import (
	"context"
	"gw-currency-wallet/config"
	"gw-currency-wallet/internal/models"
	gw_grpc "gw-currency-wallet/internal/pb/exchange"
	"gw-currency-wallet/internal/repository"
	"gw-currency-wallet/pkg"
)

type Exchange interface {
	IsExistCurrency(ctx context.Context, currency pkg.Currency) (bool, error)
	GetRates(ctx context.Context) (pkg.ExchangeRates, error)
	GetRate(ctx context.Context, from, to pkg.Currency) (pkg.Rate, error)
}

type Wallet interface {
	GetAllByEmail(ctx context.Context, email string) (pkg.AccountWallets, error)
	Deposit(ctx context.Context, email string, currency pkg.Currency, amount float32) (pkg.AccountWallets, error)
	Withdraw(ctx context.Context, email string, currency pkg.Currency, amount float32) (pkg.AccountWallets, error)
	GetRates(ctx context.Context) (pkg.ExchangeRates, error)
	Exchange(ctx context.Context, email string, from, to pkg.Currency, amount float32) (exchangedAmount float32, wallets pkg.AccountWallets, err error)
}

type Auth interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
	GenerateJWT(email string) (string, error)
	GetClaims(tokenString string) (*models.AuthClaims, error)
}

type Account interface {
	Register(ctx context.Context, email, username, password string) (*models.Account, error)
	Login(ctx context.Context, username, password string) (token string, err error)
}

type Service struct {
	Auth
	Account
	Wallet
	Exchange
}

func NewService(ctx context.Context, repo *repository.Repository, authConfig *config.AuthConfig, exchangeClient gw_grpc.ExchangeServiceClient) *Service {
	s := &Service{}

	s.Account = NewAccountService(repo.Account, s)
	s.Auth = NewAuthService(authConfig)
	s.Wallet = NewWalletService(repo.Wallet, s)
	s.Exchange = NewExchangeService(ctx, exchangeClient)

	return s
}
