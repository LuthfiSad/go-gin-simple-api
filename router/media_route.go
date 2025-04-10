package router

import (
	"go-gin-simple-api/handler"
	"go-gin-simple-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupMediaRoutes(api *gin.RouterGroup, mediaHandler *handler.MediaHandler) {
	media := api.Group("/media")
	media.GET("/", middleware.RoleAuth("admin"), mediaHandler.ListMedia)
	media.GET("/:id", middleware.RoleAuth("admin"), mediaHandler.GetMedia)
	media.POST("/", middleware.RoleAuth("admin"), mediaHandler.UploadMedia)
	media.DELETE("/:id", middleware.RoleAuth("admin"), mediaHandler.DeleteMedia)
}
