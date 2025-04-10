package router

import (
	"go-gin-simple-api/handler"
	"go-gin-simple-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(api *gin.RouterGroup, customerHandler *handler.CustomerHandler) {
	customerRoute := api.Group("/customers")
	customerRoute.GET("/", customerHandler.ListCustomers)
	customerRoute.GET("/:id", customerHandler.GetCustomer)
	customerRoute.GET("/:id/transactions", customerHandler.GetByIDWithTransactions)
	customerRoute.GET("/code/:code", customerHandler.GetByCode)

	// Protected customer routes (admin only)
	customerRoute.POST("/", middleware.RoleAuth("admin"), customerHandler.CreateCustomer)
	customerRoute.PUT("/:id", middleware.RoleAuth("admin"), customerHandler.UpdateCustomer)
	customerRoute.DELETE("/:id", middleware.RoleAuth("admin"), customerHandler.DeleteCustomer)
}
