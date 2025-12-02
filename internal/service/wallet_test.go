package service

import (
	"gw-currency-wallet/internal/db"
	gw_grpc "gw-currency-wallet/internal/pb/exchange"
	"gw-currency-wallet/internal/pb/exchange/mocks"
	mock_repository "gw-currency-wallet/internal/repository/mocks"
	"gw-currency-wallet/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDeposit_PositiveAmount_IncreasesBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGrpcExchange := mocks.NewMockExchangeServiceClient(ctrl)
	mockGrpcExchange.EXPECT().GetExchangeRates(t.Context(), nil).Return(&gw_grpc.ExchangeRatesResponse{Rates: pkg.ExchangeRates{"USD": 1}}, nil)

	mockRepo := mock_repository.NewMockWallet(ctrl)
	s := &Service{
		Exchange: NewExchangeService(t.Context(), mockGrpcExchange),
	}
	srv := NewWalletService(mockRepo, s)

	email := "user@example.com"
	currency := "USD"
	amount := float32(100)
	initialBalance := float32(0)

	mockTx := mock_repository.NewMockTx(ctrl)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil).Times(1)
	mockTx.EXPECT().Rollback(gomock.Any()).AnyTimes()

	mockRepo.EXPECT().WithTx(gomock.Any()).Return(t.Context(), mockTx, nil)

	mockRepo.EXPECT().IsExistCurrency(t.Context(), email, currency).Return(true, nil)

	mockRepo.EXPECT().GetForUpdate(t.Context(), email, currency).Return(&db.AppWallet{
		Email:    email,
		Currency: currency,
		Balance:  initialBalance,
	}, nil)

	mockRepo.EXPECT().Update(t.Context(), email, currency, initialBalance+amount).Return(nil, nil)

	mockRepo.EXPECT().GetAllByEmail(gomock.Any(), email).Return([]db.AppWallet{
		{
			Email:    email,
			Currency: currency,
			Balance:  initialBalance + amount,
		},
	}, nil)

	result, err := srv.Deposit(t.Context(), email, currency, amount)

	assert.NoError(t, err)
	assert.Equal(t, amount, result[currency])
}

func TestDeposit_ZeroAmount_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGrpcExchange := mocks.NewMockExchangeServiceClient(ctrl)
	mockGrpcExchange.EXPECT().GetExchangeRates(t.Context(), nil).Return(&gw_grpc.ExchangeRatesResponse{Rates: pkg.ExchangeRates{"USD": 1}}, nil)

	mockRepo := mock_repository.NewMockWallet(ctrl)
	s := &Service{
		Exchange: NewExchangeService(t.Context(), mockGrpcExchange),
	}
	srv := NewWalletService(mockRepo, s)

	email := "user@example.com"
	currency := "USD"
	amount := float32(0)

	mockTx := mock_repository.NewMockTx(ctrl)
	mockTx.EXPECT().Rollback(gomock.Any()).AnyTimes()

	mockRepo.EXPECT().WithTx(gomock.Any()).Return(t.Context(), mockTx, nil)

	mockRepo.EXPECT().IsExistCurrency(t.Context(), email, currency).Return(true, nil)

	mockRepo.EXPECT().GetForUpdate(t.Context(), email, currency).Return(&db.AppWallet{
		Email:    email,
		Currency: currency,
		Balance:  0,
	}, nil)

	_, err := srv.Deposit(t.Context(), email, currency, amount)

	assert.ErrorIs(t, err, ErrZeroAmount)
}

func TestDeposit_NegativeAmount_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGrpcExchange := mocks.NewMockExchangeServiceClient(ctrl)
	mockGrpcExchange.EXPECT().GetExchangeRates(t.Context(), nil).Return(&gw_grpc.ExchangeRatesResponse{Rates: pkg.ExchangeRates{"USD": 1}}, nil)

	mockRepo := mock_repository.NewMockWallet(ctrl)
	s := &Service{
		Exchange: NewExchangeService(t.Context(), mockGrpcExchange),
	}
	srv := NewWalletService(mockRepo, s)

	email := "user@example.com"
	currency := "USD"
	amount := float32(-1)

	mockTx := mock_repository.NewMockTx(ctrl)
	mockTx.EXPECT().Rollback(gomock.Any()).AnyTimes()

	mockRepo.EXPECT().WithTx(gomock.Any()).Return(t.Context(), mockTx, nil)

	mockRepo.EXPECT().IsExistCurrency(t.Context(), email, currency).Return(true, nil)

	mockRepo.EXPECT().GetForUpdate(t.Context(), email, currency).Return(&db.AppWallet{
		Email:    email,
		Currency: currency,
		Balance:  0,
	}, nil)

	_, err := srv.Deposit(t.Context(), email, currency, amount)

	assert.ErrorIs(t, err, ErrNegativeAmount)
}

func TestWithdraw_PositiveAmount_DecreasesBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGrpcExchange := mocks.NewMockExchangeServiceClient(ctrl)
	mockGrpcExchange.EXPECT().GetExchangeRates(t.Context(), nil).Return(&gw_grpc.ExchangeRatesResponse{Rates: pkg.ExchangeRates{"USD": 1}}, nil)

	mockRepo := mock_repository.NewMockWallet(ctrl)
	s := &Service{
		Exchange: NewExchangeService(t.Context(), mockGrpcExchange),
	}
	srv := NewWalletService(mockRepo, s)

	email := "user@example.com"
	currency := "USD"
	amount := float32(20)
	initialBalance := float32(100)

	mockTx := mock_repository.NewMockTx(ctrl)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil).Times(1)
	mockTx.EXPECT().Rollback(gomock.Any()).AnyTimes()

	mockRepo.EXPECT().WithTx(gomock.Any()).Return(t.Context(), mockTx, nil)

	mockRepo.EXPECT().IsExistCurrency(t.Context(), email, currency).Return(true, nil)

	mockRepo.EXPECT().GetForUpdate(t.Context(), email, currency).Return(&db.AppWallet{
		Email:    email,
		Currency: currency,
		Balance:  initialBalance,
	}, nil)

	mockRepo.EXPECT().Update(t.Context(), email, currency, initialBalance-amount).Return(nil, nil)

	mockRepo.EXPECT().GetAllByEmail(gomock.Any(), email).Return([]db.AppWallet{
		{
			Email:    email,
			Currency: currency,
			Balance:  initialBalance - amount,
		},
	}, nil)

	result, err := srv.Withdraw(t.Context(), email, currency, amount)

	assert.NoError(t, err)
	assert.Equal(t, initialBalance-amount, result[currency])
}

func TestWithdraw_ZeroAmount_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGrpcExchange := mocks.NewMockExchangeServiceClient(ctrl)
	mockGrpcExchange.EXPECT().GetExchangeRates(t.Context(), nil).Return(&gw_grpc.ExchangeRatesResponse{Rates: pkg.ExchangeRates{"USD": 1}}, nil)

	mockRepo := mock_repository.NewMockWallet(ctrl)
	s := &Service{
		Exchange: NewExchangeService(t.Context(), mockGrpcExchange),
	}
	srv := NewWalletService(mockRepo, s)

	email := "user@example.com"
	currency := "USD"
	amount := float32(0)

	mockTx := mock_repository.NewMockTx(ctrl)
	mockTx.EXPECT().Rollback(gomock.Any()).AnyTimes()

	mockRepo.EXPECT().WithTx(gomock.Any()).Return(t.Context(), mockTx, nil)

	mockRepo.EXPECT().IsExistCurrency(t.Context(), email, currency).Return(true, nil)

	mockRepo.EXPECT().GetForUpdate(t.Context(), email, currency).Return(&db.AppWallet{
		Email:    email,
		Currency: currency,
		Balance:  0,
	}, nil)

	_, err := srv.Withdraw(t.Context(), email, currency, amount)

	assert.ErrorIs(t, err, ErrZeroAmount)
}

func TestWithdraw_NegativeAmount_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGrpcExchange := mocks.NewMockExchangeServiceClient(ctrl)
	mockGrpcExchange.EXPECT().GetExchangeRates(t.Context(), nil).Return(&gw_grpc.ExchangeRatesResponse{Rates: pkg.ExchangeRates{"USD": 1}}, nil)

	mockRepo := mock_repository.NewMockWallet(ctrl)
	s := &Service{
		Exchange: NewExchangeService(t.Context(), mockGrpcExchange),
	}
	srv := NewWalletService(mockRepo, s)

	email := "user@example.com"
	currency := "USD"
	amount := float32(-1)

	mockTx := mock_repository.NewMockTx(ctrl)
	mockTx.EXPECT().Rollback(gomock.Any()).AnyTimes()

	mockRepo.EXPECT().WithTx(gomock.Any()).Return(t.Context(), mockTx, nil)

	mockRepo.EXPECT().IsExistCurrency(t.Context(), email, currency).Return(true, nil)

	mockRepo.EXPECT().GetForUpdate(t.Context(), email, currency).Return(&db.AppWallet{
		Email:    email,
		Currency: currency,
		Balance:  0,
	}, nil)

	_, err := srv.Withdraw(t.Context(), email, currency, amount)

	assert.ErrorIs(t, err, ErrNegativeAmount)
}

func TestWithdraw_InsufficientBalance_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGrpcExchange := mocks.NewMockExchangeServiceClient(ctrl)
	mockGrpcExchange.EXPECT().GetExchangeRates(t.Context(), nil).Return(&gw_grpc.ExchangeRatesResponse{Rates: pkg.ExchangeRates{"USD": 1}}, nil)

	mockRepo := mock_repository.NewMockWallet(ctrl)
	s := &Service{
		Exchange: NewExchangeService(t.Context(), mockGrpcExchange),
	}
	srv := NewWalletService(mockRepo, s)

	email := "user@example.com"
	currency := "USD"
	amount := float32(101)
	initialBalance := float32(100)

	mockTx := mock_repository.NewMockTx(ctrl)
	mockTx.EXPECT().Rollback(gomock.Any()).AnyTimes()

	mockRepo.EXPECT().WithTx(gomock.Any()).Return(t.Context(), mockTx, nil)

	mockRepo.EXPECT().IsExistCurrency(t.Context(), email, currency).Return(true, nil)

	mockRepo.EXPECT().GetForUpdate(t.Context(), email, currency).Return(&db.AppWallet{
		Email:    email,
		Currency: currency,
		Balance:  initialBalance,
	}, nil)

	_, err := srv.Withdraw(t.Context(), email, currency, amount)

	assert.ErrorIs(t, err, ErrInsufficientBalance)
}
