package router

import (
	"go-gin-simple-api/handler"
	"go-gin-simple-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupBookStockRoutes(api *gin.RouterGroup, bookStockHandler *handler.BookStockHandler) {
	bookStock := api.Group("/bookstocks")
	bookStock.GET("/", bookStockHandler.ListBookStocks)
	bookStock.GET("/:code", bookStockHandler.GetBookStock)
	bookStock.GET("/book/:book_id", bookStockHandler.GetByBookID)
	bookStock.GET("/book/:book_id/available", bookStockHandler.GetAvailableByBookID)

	// Protected routes (admin only)
	bookStock.POST("/", middleware.RoleAuth("admin"), bookStockHandler.CreateBookStock)
	bookStock.PUT("/:code", middleware.RoleAuth("admin"), bookStockHandler.UpdateBookStock)
	bookStock.DELETE("/:code", middleware.RoleAuth("admin"), bookStockHandler.DeleteBookStock)
	bookStock.PATCH("/:code/status", middleware.RoleAuth("admin"), bookStockHandler.UpdateStatus)
}
