package handler

import (
	_ "gw-currency-wallet/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (h *Handler) Router() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("api/v1")
	{
		v1.POST("register", h.Register)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
