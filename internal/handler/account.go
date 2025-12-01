package handler

import (
	"errors"
	"gw-currency-wallet/internal/dto"
	"gw-currency-wallet/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Register godoc
// @Summary Регистрация пользователя
// @Description Регистрация нового пользователя.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration data"
// @Success 201 {object} dto.Message "User registered successfully"
// @Failure 400 {object} dto.Message "Username or email already exists"
// @Router /api/v1/register [post]
func (h *Handler) Register(c *gin.Context) {
	var in dto.RegisterRequest

	if err := c.BindJSON(&in); err != nil {
		sendBadRequest(c, err)
		return
	}

	_, err := h.s.Account.Register(c, in.Email, in.Username, in.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailAlreadyExists):
			sendBadRequest(c, err)
		case errors.Is(err, service.ErrUsernameAlreadyExists):
			sendBadRequest(c, err)
		default:
			zap.L().Error(err.Error())
			sendInternalError(c)
		}
		return
	}

	sendCreated(c, "User registered successfully")
}
