package handler

import (
	"errors"
	"gw-currency-wallet/internal/dto"
	"gw-currency-wallet/internal/service"

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

// Deposit godoc
// @Summary Пополнение счета пользователя
// @Description Позволяет пользователю пополнить свой счет. Проверяется корректность суммы и валюты.
// @Tags wallet
// @Accept json
// @Produce json
// @Param input body dto.DepositRequest true "Сумма и валюта для пополнения"
// @Success 200 {object} dto.DepositResponse "Account topped up successfully"
// @Failure 400 {object} dto.Message "Invalid amount or currency"
// @Failure 500 {object} dto.Message "Internal server error"
// @Router /api/v1/wallet/deposit [post]
// @Security BearerAuth
func (h *Handler) Deposit(c *gin.Context) {
	var in dto.DepositRequest

	if err := c.BindJSON(&in); err != nil {
		sendBadRequest(c, err)
		return
	}

	email, ok := getAccountFromContext(c)
	if !ok {
		sendInternalError(c)
		return
	}

	wallets, err := h.s.Deposit(c, email, in.Currency, in.Amount)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNegativeAmount):
			sendBadRequest(c, service.ErrNegativeAmount)
			return
		case errors.Is(err, service.ErrZeroAmount):
			sendBadRequest(c, service.ErrZeroAmount)
			return
		case errors.Is(err, service.ErrNonExistentCurrency):
			sendBadRequest(c, service.ErrNonExistentCurrency)
			return
		default:
			zap.L().Error(err.Error())
			sendInternalError(c)
			return
		}
	}

	sendOK(c, &dto.DepositResponse{
		Message:    "Account topped up successfully",
		NewBalance: wallets,
	})
}

// Withdraw godoc
// @Summary Вывод средств со счета пользователя
// @Description Позволяет пользователю вывести средства со своего счета.
// @Tags wallet
// @Accept json
// @Produce json
// @Param input body dto.WithdrawRequest true "Сумма и валюта для вывода"
// @Success 200 {object} dto.WithdrawResponse "Withdrawal successful"
// @Failure 400 {object} dto.Message "Insufficient funds or invalid amount"
// @Failure 500 {object} dto.Message "Internal server error"
// @Router /api/v1/wallet/withdraw [post]
// @Security BearerAuth
func (h *Handler) Withdraw(c *gin.Context) {
	var in dto.WithdrawRequest

	if err := c.BindJSON(&in); err != nil {
		sendBadRequest(c, err)
		return
	}

	email, ok := getAccountFromContext(c)
	if !ok {
		sendInternalError(c)
		return
	}

	wallets, err := h.s.Withdraw(c, email, in.Currency, in.Amount)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNegativeAmount):
			sendBadRequest(c, service.ErrNegativeAmount)
			return
		case errors.Is(err, service.ErrZeroAmount):
			sendBadRequest(c, service.ErrZeroAmount)
			return
		case errors.Is(err, service.ErrNonExistentCurrency):
			sendBadRequest(c, service.ErrNonExistentCurrency)
			return
		case errors.Is(err, service.ErrInsufficientBalance):
			sendBadRequest(c, service.ErrInsufficientBalance)
			return
		default:
			zap.L().Error(err.Error())
			sendInternalError(c)
			return
		}
	}

	sendOK(c, &dto.WithdrawResponse{
		Message:    "Withdrawal successful",
		NewBalance: wallets,
	})
}

// GetRates godoc
// @Summary Получение актуальных курсов валют
// @Description Возвращает курсы всех поддерживаемых валют.
// @Tags exchange
// @Accept json
// @Produce json
// @Success 200 {object} dto.GetRatesResponse "Exchange rates"
// @Failure 500 {object} dto.Message "Failed to retrieve exchange rates"
// @Router /api/v1/exchange/rates [get]
// @Security BearerAuth
func (h *Handler) GetRates(c *gin.Context) {
	rates, err := h.s.Wallet.GetRates(c)
	if err != nil {
		sendInternalError(c)
		return
	}

	sendOK(c, &dto.GetRatesResponse{
		Rates: rates,
	})
}
