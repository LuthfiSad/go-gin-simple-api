package router

import (
	"go-gin-simple-api/handler"
	"go-gin-simple-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupBookTransactionRoutes(api *gin.RouterGroup, bookTransactionHandler *handler.BookTransactionHandler) {
	transactionRoute := api.Group("/transactions")
	transactionRoute.GET("", bookTransactionHandler.ListBookTransactions)
	transactionRoute.GET("/:id", bookTransactionHandler.GetBookTransaction)
	transactionRoute.GET("/customer/:customer_id", bookTransactionHandler.GetByCustomerID)
	transactionRoute.GET("/book/:book_id", bookTransactionHandler.GetByBookID)
	transactionRoute.GET("/stock/:stock_code", bookTransactionHandler.GetByStockCode)
	transactionRoute.GET("/overdue", bookTransactionHandler.GetOverdueTransactions)

	// Protected book transaction routes (admin only)
	transactionRoute.POST("", middleware.RoleAuth("admin"), bookTransactionHandler.CreateBookTransaction)
	transactionRoute.PUT("/:id", middleware.RoleAuth("admin"), bookTransactionHandler.UpdateBookTransaction)
	transactionRoute.DELETE("/:id", middleware.RoleAuth("admin"), bookTransactionHandler.DeleteBookTransaction)
	transactionRoute.PATCH("/:id/status", middleware.RoleAuth("admin"), bookTransactionHandler.UpdateStatus)
	transactionRoute.POST("/:id/return", middleware.RoleAuth("admin"), bookTransactionHandler.ReturnBook)
}
