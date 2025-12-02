package handler

import (
	"gw-currency-wallet/internal/dto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetWallets godoc
// @Summary Получение кошельков пользователя
// @Description Возвращает список кошельков и их балансы для авторизованного пользователя.
// @Tags wallet
// @Accept json
// @Produce json
// @Success 200 {object} dto.GetWalletsResponse "User wallets"
// @Failure 500 {object} dto.Message "Internal server error"
// @Router /api/v1/balance [get]
// @Security BearerAuth
func (h *Handler) GetWallets(c *gin.Context) {
	email, ok := getAccountFromContext(c)
	if !ok {
		sendInternalError(c)
		return
	}

	wallets, err := h.s.Wallet.GetAllByEmail(c, email)
	if err != nil {
		zap.L().Error(err.Error())
		sendInternalError(c)
		return
	}

	sendOK(c, &dto.GetWalletsResponse{
		Balance: wallets,
	})
}
