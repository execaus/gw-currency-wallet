package service

import (
	"context"
	"errors"
	"gw-currency-wallet/internal/repository"
	"gw-currency-wallet/pkg"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type WalletService struct {
	r repository.Wallet
	s *Service
}

func (s *WalletService) Exchange(ctx context.Context, email string, from, to pkg.Currency, amount float32) (exchangedAmount float32, wallets pkg.AccountWallets, err error) {
	c, tx, err := s.r.WithTx(ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return 0, nil, err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			zap.L().Error(err.Error())
		}
	}()

	rate, err := s.s.Exchange.GetRate(c, from, to)
	if err != nil {
		zap.L().Error(err.Error())
		return 0, nil, err
	}

	exchangedAmount = amount * rate

	if err = s.withdraw(c, email, from, amount); err != nil {
		zap.L().Error(err.Error())
		return 0, nil, err
	}

	if err = s.deposit(c, email, to, exchangedAmount); err != nil {
		zap.L().Error(err.Error())
		return 0, nil, err
	}

	if err = tx.Commit(c); err != nil {
		zap.L().Error(err.Error())
		return 0, nil, err
	}

	wallets, err = s.accountWallets(ctx, email)

	return exchangedAmount, wallets, err
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
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			zap.L().Error(err.Error())
		}
	}()

	if err = s.withdraw(c, email, currency, amount); err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	if err = tx.Commit(c); err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	return s.accountWallets(ctx, email)
}

func (s *WalletService) Deposit(ctx context.Context, email string, currency pkg.Currency, amount float32) (pkg.AccountWallets, error) {
	c, tx, err := s.r.WithTx(ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			zap.L().Error(err.Error())
		}
	}()

	if err = s.deposit(c, email, currency, amount); err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	if err = tx.Commit(c); err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	return s.accountWallets(ctx, email)
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

func (s *WalletService) deposit(ctx context.Context, email string, currency pkg.Currency, amount float32) error {
	isExistsCurrency, err := s.s.Exchange.IsExistCurrency(ctx, currency)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	if !isExistsCurrency {
		zap.L().Error(ErrNonExistentCurrency.Error())
		return ErrNonExistentCurrency
	}

	isExistWallet, err := s.r.IsExistCurrency(ctx, email, currency)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	if !isExistWallet {
		if err = s.r.Create(ctx, email, currency); err != nil {
			zap.L().Error(err.Error())
			return err
		}
	}

	wallet, err := s.r.GetForUpdate(ctx, email, currency)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	if amount == 0 {
		zap.L().Error(ErrZeroAmount.Error())
		return ErrZeroAmount
	}
	if amount < 0 {
		zap.L().Error(ErrNegativeAmount.Error())
		return ErrNegativeAmount
	}

	newBalance := wallet.Balance + amount

	_, err = s.r.Update(ctx, email, currency, newBalance)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	return nil
}

func (s *WalletService) withdraw(ctx context.Context, email string, currency pkg.Currency, amount float32) error {
	isExistsCurrency, err := s.s.Exchange.IsExistCurrency(ctx, currency)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	if !isExistsCurrency {
		zap.L().Error(ErrNonExistentCurrency.Error())
		return ErrNonExistentCurrency
	}

	isExistWallet, err := s.r.IsExistCurrency(ctx, email, currency)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	if !isExistWallet {
		if err = s.r.Create(ctx, email, currency); err != nil {
			zap.L().Error(err.Error())
			return err
		}
	}

	wallet, err := s.r.GetForUpdate(ctx, email, currency)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	if amount == 0 {
		zap.L().Error(ErrZeroAmount.Error())
		return ErrZeroAmount
	}
	if amount < 0 {
		zap.L().Error(ErrNegativeAmount.Error())
		return ErrNegativeAmount
	}
	if wallet.Balance < amount {
		zap.L().Error(ErrInsufficientBalance.Error())
		return ErrInsufficientBalance
	}

	newBalance := wallet.Balance - amount

	_, err = s.r.Update(ctx, email, currency, newBalance)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	return nil
}

func (s *WalletService) accountWallets(ctx context.Context, email string) (pkg.AccountWallets, error) {
	wallets, err := s.r.GetAllByEmail(ctx, email)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	result := make(pkg.AccountWallets, len(wallets))
	for _, w := range wallets {
		result[w.Currency] = w.Balance
	}

	return result, nil
}
