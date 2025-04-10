package router

import (
	"go-gin-simple-api/handler"
	"go-gin-simple-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupBookRoutes(api *gin.RouterGroup, bookHandler *handler.BookHandler) {
	bookRoute := api.Group("/books")
	bookRoute.GET("/", bookHandler.ListBooks)
	bookRoute.GET("/:id", bookHandler.GetBook)
	bookRoute.POST("/", middleware.RoleAuth("admin"), bookHandler.CreateBook)
	bookRoute.PUT("/:id", middleware.RoleAuth("admin"), bookHandler.UpdateBook)
	bookRoute.DELETE("/:id", middleware.RoleAuth("admin"), bookHandler.DeleteBook)
	bookRoute.DELETE("/:id/cover", middleware.RoleAuth("admin"), bookHandler.DeleteBookCover)
}
