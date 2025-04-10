package router

import (
	"go-gin-simple-api/handler"
	"go-gin-simple-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupChargeRoutes(api *gin.RouterGroup, chargeHandler *handler.ChargeHandler) {
	chargeRoute := api.Group("/charges")
	chargeRoute.GET("", chargeHandler.ListCharges)
	chargeRoute.GET("/:id", chargeHandler.GetCharge)
	chargeRoute.GET("/transaction/:transaction_id", chargeHandler.GetByBookTransactionID)
	chargeRoute.GET("/user/:user_id", chargeHandler.GetByUserID)

	// Protected charge routes
	chargeRoute.POST("/", middleware.RoleAuth("admin"), chargeHandler.CreateCharge) // This already checks userData
	chargeRoute.PUT("/:id", middleware.RoleAuth("admin"), chargeHandler.UpdateCharge)
	chargeRoute.DELETE("/:id", middleware.RoleAuth("admin"), chargeHandler.DeleteCharge)
}
