package router

import (
	"go-gin-simple-api/handler"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(api *gin.RouterGroup, authHandler *handler.AuthHandler) {
	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)
}
