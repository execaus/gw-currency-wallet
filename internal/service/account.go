package service

import (
	"context"
	"gw-currency-wallet/internal/models"
	"gw-currency-wallet/internal/repository"

	"go.uber.org/zap"
)

type AccountService struct {
	r repository.Account
	s *Service
}

func (s *AccountService) Login(ctx context.Context, username, password string) (string, error) {
	account, err := s.r.GetByUsername(ctx, username)
	if err != nil {
		zap.L().Error(err.Error())
		return "", err
	}

	if account == nil {
		zap.L().Warn("invalid username")
		return "", ErrInvalidCredentials
	}

	if err = s.s.Auth.ComparePassword(account.Password, password); err != nil {
		zap.L().Warn("invalid password")
		return "", ErrInvalidCredentials
	}

	token, err := s.s.Auth.GenerateJWT(account.Email)
	if err != nil {
		zap.L().Error(err.Error())
		return "", err
	}

	return token, nil
}

func (s *AccountService) Register(ctx context.Context, email, username, password string) (*models.Account, error) {
	passwordHash, err := s.s.Auth.HashPassword(password)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	isExist, err := s.r.IsUsernameExist(ctx, username)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	if isExist {
		zap.L().Warn(ErrUsernameAlreadyExists.Error())
		return nil, ErrUsernameAlreadyExists
	}

	isExist, err = s.r.IsEmailExist(ctx, username)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	if isExist {
		zap.L().Warn(ErrEmailAlreadyExists.Error())
		return nil, ErrEmailAlreadyExists
	}

	dbAccount, err := s.r.Create(ctx, email, username, passwordHash)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	return &models.Account{
		Email:        dbAccount.Email,
		PasswordHash: passwordHash,
		Username:     dbAccount.Username,
	}, err
}

func NewAccountService(repo repository.Account, srv *Service) *AccountService {
	return &AccountService{
		s: srv,
		r: repo,
	}
}
