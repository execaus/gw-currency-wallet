package dto

import "gw-currency-wallet/internal/models"

type GetWalletsResponse struct {
	Balance models.AccountWallets `json:"balance"`
}
