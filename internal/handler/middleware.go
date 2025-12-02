package handler

import (
	"errors"
	"gw-currency-wallet/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type contextKey string

const AccountEmailKey contextKey = "accountEmail"

var (
	ErrInvalidAuthorizationHeader = errors.New("missing or invalid Authorization header")
)

func (h *Handler) authMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		zap.L().Error(ErrInvalidAuthorizationHeader.Error())
		sendUnauthorized(c, ErrInvalidAuthorizationHeader)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := h.s.Auth.GetClaims(token)
	if err != nil {
		zap.L().Error(service.ErrTokenInvalid.Error())
		sendUnauthorized(c, service.ErrTokenInvalid)
		return
	}

	c.Set(AccountEmailKey, claims.Email)
	c.Next()
}

func getAccountFromContext(ctx *gin.Context) (string, bool) {
	accountID, ok := ctx.Get(AccountEmailKey)
	if !ok {
		return "", false
	}
	idStr, ok := accountID.(string)
	return idStr, ok
}
