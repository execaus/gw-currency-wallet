package service

import (
	"context"
	"gw-currency-wallet/config"
	"gw-currency-wallet/internal/models"
	"gw-currency-wallet/internal/repository"
)

type Auth interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
	GenerateJWT(userID string) (string, error)
	GetClaims(tokenString string) (*models.AuthClaims, error)
}

type Account interface {
	Register(ctx context.Context, email, username, password string) (*models.Account, error)
}

type Service struct {
	Auth
	Account
}

func NewService(repo repository.Repository, authConfig *config.AuthConfig) *Service {
	s := Service{}

	s.Account = NewAccountService(repo.Account, &s)
	s.Auth = NewAuthService(authConfig)

	return &s
}
