package handler

import (
	"gw-currency-wallet/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func sendBadRequest(c *gin.Context, err error) {
	send(c, http.StatusBadRequest, dto.ErrorMessage{Error: err.Error()})
}

func sendCreated(c *gin.Context, message string) {
	send(c, http.StatusCreated, dto.Message{Message: message})
}

func sendInternalError(c *gin.Context) {
	send(c, http.StatusInternalServerError, dto.Message{Message: "server error"})
}

func sendOK(c *gin.Context, body any) {
	send(c, http.StatusOK, body)
}

func sendUnauthorized(c *gin.Context, err error) {
	send(c, http.StatusUnauthorized, dto.Message{Message: err.Error()})
}

func send(c *gin.Context, status int, body any) {
	zap.L().Info("HTTP Response", zap.Int("status", status), zap.Any("body", body))
	c.AbortWithStatusJSON(status, body)
}
