package service

import (
	"context"
	"gw-currency-wallet/internal/repository"
	"gw-currency-wallet/pkg"

	"go.uber.org/zap"
)

type WalletService struct {
	r repository.Wallet
	s *Service
}

func (s *WalletService) GetRates(ctx context.Context) (pkg.ExchangeRates, error) {
	rates, err := s.s.Exchange.GetRates(ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	return rates, err
}

func (s *WalletService) Withdraw(ctx context.Context, email string, currency pkg.Currency, amount float32) (pkg.AccountWallets, error) {
	c, tx, err := s.r.WithTx(ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			zap.L().Error(err.Error())
		}
	}()

	isExistsCurrency, err := s.s.Exchange.IsExistCurrency(ctx, currency)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	if !isExistsCurrency {
		zap.L().Error(ErrNonExistentCurrency.Error())
		return nil, ErrNonExistentCurrency
	}

	isExistWallet, err := s.r.IsExistCurrency(c, email, currency)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	if !isExistWallet {
		if err = s.r.Create(c, email, currency); err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
	}

	wallet, err := s.r.GetForUpdate(c, email, currency)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	if amount == 0 {
		zap.L().Error(ErrZeroAmount.Error())
		return nil, ErrZeroAmount
	}
	if amount < 0 {
		zap.L().Error(ErrNegativeAmount.Error())
		return nil, ErrNegativeAmount
	}
	if wallet.Balance < amount {
		zap.L().Error(ErrInsufficientBalance.Error())
		return nil, ErrInsufficientBalance
	}

	newBalance := wallet.Balance - amount

	_, err = s.r.Update(c, email, currency, newBalance)
	if err != nil {
		zap.L().Error(ErrNegativeAmount.Error())
		return nil, err
	}

	if err = tx.Commit(c); err != nil {
		zap.L().Error(ErrNegativeAmount.Error())
		return nil, err
	}

	wallets, err := s.r.GetAllByEmail(ctx, email)
	if err != nil {
		zap.L().Error(ErrNegativeAmount.Error())
		return nil, err
	}

	result := make(pkg.AccountWallets, len(wallets))
	for _, w := range wallets {
		result[w.Currency] = w.Balance
	}

	return result, nil
}

func (s *WalletService) Deposit(ctx context.Context, email string, currency pkg.Currency, amount float32) (pkg.AccountWallets, error) {
	c, tx, err := s.r.WithTx(ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			zap.L().Error(err.Error())
		}
	}()

	isExistsCurrency, err := s.s.Exchange.IsExistCurrency(ctx, currency)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	if !isExistsCurrency {
		zap.L().Error(ErrNonExistentCurrency.Error())
		return nil, ErrNonExistentCurrency
	}

	isExistWallet, err := s.r.IsExistCurrency(c, email, currency)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	if !isExistWallet {
		if err = s.r.Create(c, email, currency); err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
	}

	wallet, err := s.r.GetForUpdate(c, email, currency)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	if amount == 0 {
		zap.L().Error(ErrZeroAmount.Error())
		return nil, ErrZeroAmount
	}
	if amount < 0 {
		zap.L().Error(ErrNegativeAmount.Error())
		return nil, ErrNegativeAmount
	}

	newBalance := wallet.Balance + amount

	_, err = s.r.Update(c, email, currency, newBalance)
	if err != nil {
		zap.L().Error(ErrNegativeAmount.Error())
		return nil, err
	}

	if err = tx.Commit(c); err != nil {
		zap.L().Error(ErrNegativeAmount.Error())
		return nil, err
	}

	wallets, err := s.r.GetAllByEmail(ctx, email)
	if err != nil {
		zap.L().Error(ErrNegativeAmount.Error())
		return nil, err
	}

	result := make(pkg.AccountWallets, len(wallets))
	for _, w := range wallets {
		result[w.Currency] = w.Balance
	}

	return result, nil
}

func (s *WalletService) GetAllByEmail(ctx context.Context, email string) (pkg.AccountWallets, error) {
	wallets, err := s.r.GetAllByEmail(ctx, email)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	if wallets == nil {
		return make(pkg.AccountWallets), nil
	}

	result := make(pkg.AccountWallets, len(wallets))
	for _, wallet := range wallets {
		result[wallet.Currency] = wallet.Balance
	}

	return result, nil
}

func NewWalletService(r repository.Wallet, s *Service) *WalletService {
	return &WalletService{
		r: r,
		s: s,
	}
}
