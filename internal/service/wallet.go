package service

import (
	"context"
	"gw-currency-wallet/internal/models"
	"gw-currency-wallet/internal/repository"

	"go.uber.org/zap"
)

type WalletService struct {
	r repository.Wallet
}

func (s *WalletService) GetAllByEmail(ctx context.Context, email string) (models.AccountWallets, error) {
	wallets, err := s.r.GetAllByEmail(ctx, email)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	if wallets == nil {
		return make(models.AccountWallets), nil
	}

	result := make(models.AccountWallets, len(wallets))
	for _, wallet := range wallets {
		result[wallet.Currency] = wallet.Balance
	}

	return result, nil
}

func NewWalletService(r repository.Wallet) *WalletService {
	return &WalletService{
		r: r,
	}
}
