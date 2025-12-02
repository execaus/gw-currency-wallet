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
