package dto

import (
	"gw-currency-wallet/pkg"
)

type GetWalletsResponse struct {
	Balance pkg.AccountWallets `json:"balance"`
}

type DepositRequest struct {
	Amount   float32 `json:"amount" binding:"required,gt=0"`
	Currency string  `json:"currency" binding:"required"`
}

type DepositResponse struct {
	Message    string             `json:"message"`
	NewBalance pkg.AccountWallets `json:"new_balance"`
}

type WithdrawRequest struct {
	Amount   float32 `json:"amount" binding:"required,gt=0"`
	Currency string  `json:"currency" binding:"required"`
}

type WithdrawResponse struct {
	Message    string             `json:"message"`
	NewBalance pkg.AccountWallets `json:"new_balance"`
}

type GetRatesResponse struct {
	Rates pkg.ExchangeRates `json:"rates"`
}

type ExchangeRequest struct {
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Amount       float32 `json:"amount"`
}

type ExchangeResponse struct {
	Message         string             `json:"message"`
	ExchangedAmount float32            `json:"exchanged_amount"`
	NewBalance      pkg.AccountWallets `json:"new_balance"`
}
