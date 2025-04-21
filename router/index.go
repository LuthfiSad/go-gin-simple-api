package router

import (
	"go-gin-simple-api/handler"
	"go-gin-simple-api/lib"
	"go-gin-simple-api/repository"
	"go-gin-simple-api/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cloudinary *lib.CloudinaryService) *gin.Engine {
	router := gin.Default()
	repo := repository.NewRepository(db)
	service := service.NewService(repo, cloudinary)
	// menggunakan lazzy
	handlerInstance := handler.NewHandler(service)
	// handler := handler.NewHandler(service)

	api := router.Group("/api")
	SetupAuthRoutes(api, handlerInstance.GetAuthHandler())
	// SetupAuthRoutes(api, &handler.AuthHandler)

	// api.POST("/register", handler.AuthHandler.Register)
	// api.POST("/login", handler.AuthHandler.Login)

	// api.Use(middleware.JWTAuth(repo.AuthRepository))
	// menggunakan lazzy
	SetupBookRoutes(api, handlerInstance.GetBookHandler())
	SetupMediaRoutes(api, handlerInstance.GetMediaHandler())
	SetupBookStockRoutes(api, handlerInstance.GetBookStockHandler())
	SetupCustomerRoutes(api, handlerInstance.GetCustomerHandler())
	SetupChargeRoutes(api, handlerInstance.GetChargeHandler())
	SetupBookTransactionRoutes(api, handlerInstance.GetBookTransactionHandler())

	// SetupBookRoutes(api, &handler.BookHandler)
	// SetupMediaRoutes(api, &handler.MediaHandler)
	// SetupBookStockRoutes(api, &handler.BookStockHandler)
	// SetupCustomerRoutes(api, &handler.CustomerHandler)
	// SetupChargeRoutes(api, &handler.ChargeHandler)
	// SetupBookTransactionRoutes(api, &handler.BookTransactionHandler)

	// // Book routes
	// bookRoute := api.Group("/books")
	// bookRoute.GET("/", handler.BookHandler.ListBooks)
	// bookRoute.GET("/:id", handler.BookHandler.GetBook)
	// bookRoute.POST("/", middleware.RoleAuth("admin"), handler.BookHandler.CreateBook)
	// bookRoute.PUT("/:id", middleware.RoleAuth("admin"), handler.BookHandler.UpdateBook)
	// bookRoute.DELETE("/:id", middleware.RoleAuth("admin"), handler.BookHandler.DeleteBook)
	// bookRoute.DELETE("/:id/cover", middleware.RoleAuth("admin"), handler.BookHandler.DeleteBookCover)

	// // Media routes
	// media := api.Group("/media")
	// media.GET("/", middleware.RoleAuth("admin"), handler.MediaHandler.ListMedia)
	// media.GET("/:id", middleware.RoleAuth("admin"), handler.MediaHandler.GetMedia)
	// media.POST("/", middleware.RoleAuth("admin"), handler.MediaHandler.UploadMedia)
	// media.DELETE("/:id", middleware.RoleAuth("admin"), handler.MediaHandler.DeleteMedia)

	// // BookStock routes
	// bookStock := api.Group("/bookstocks")
	// bookStock.GET("", handler.BookStockHandler.ListBookStocks)
	// bookStock.GET("/:code", handler.BookStockHandler.GetBookStock)
	// bookStock.GET("/book/:book_id", handler.BookStockHandler.GetByBookID)
	// bookStock.GET("/book/:book_id/available", handler.BookStockHandler.GetAvailableByBookID)

	// // Protected routes (admin only)
	// bookStock.POST("", middleware.RoleAuth("admin"), handler.BookStockHandler.CreateBookStock)
	// bookStock.PUT("/:code", middleware.RoleAuth("admin"), handler.BookStockHandler.UpdateBookStock)
	// bookStock.DELETE("/:code", middleware.RoleAuth("admin"), handler.BookStockHandler.DeleteBookStock)
	// bookStock.PATCH("/:code/status", middleware.RoleAuth("admin"), handler.BookStockHandler.UpdateStatus)

	// // Customer routes
	// customerRoute := api.Group("/customers")
	// customerRoute.GET("", handler.CustomerHandler.ListCustomers)
	// customerRoute.GET("/:id", handler.CustomerHandler.GetCustomer)
	// customerRoute.GET("/:id/transactions", handler.CustomerHandler.GetByIDWithTransactions)
	// customerRoute.GET("/code/:code", handler.CustomerHandler.GetByCode)

	// // Protected customer routes (admin only)
	// customerRoute.POST("", middleware.RoleAuth("admin"), handler.CustomerHandler.CreateCustomer)
	// customerRoute.PUT("/:id", middleware.RoleAuth("admin"), handler.CustomerHandler.UpdateCustomer)
	// customerRoute.DELETE("/:id", middleware.RoleAuth("admin"), handler.CustomerHandler.DeleteCustomer)

	// // Charge routes
	// chargeRoute := api.Group("/charges")
	// chargeRoute.GET("", handler.ChargeHandler.ListCharges)
	// chargeRoute.GET("/:id", handler.ChargeHandler.GetCharge)
	// chargeRoute.GET("/transaction/:transaction_id", handler.ChargeHandler.GetByBookTransactionID)
	// chargeRoute.GET("/user/:user_id", handler.ChargeHandler.GetByUserID)

	// // Protected charge routes
	// chargeRoute.POST("", handler.ChargeHandler.CreateCharge) // This already checks userData
	// chargeRoute.PUT("/:id", middleware.RoleAuth("admin"), handler.ChargeHandler.UpdateCharge)
	// chargeRoute.DELETE("/:id", middleware.RoleAuth("admin"), handler.ChargeHandler.DeleteCharge)

	// // Book transaction routes
	// transactionRoute := api.Group("/transactions")
	// transactionRoute.GET("", handler.BookTransactionHandler.ListBookTransactions)
	// transactionRoute.GET("/:id", handler.BookTransactionHandler.GetBookTransaction)
	// transactionRoute.GET("/customer/:customer_id", handler.BookTransactionHandler.GetByCustomerID)
	// transactionRoute.GET("/book/:book_id", handler.BookTransactionHandler.GetByBookID)
	// transactionRoute.GET("/stock/:stock_code", handler.BookTransactionHandler.GetByStockCode)
	// transactionRoute.GET("/overdue", handler.BookTransactionHandler.GetOverdueTransactions)

	// // Protected book transaction routes (admin only)
	// transactionRoute.POST("", middleware.RoleAuth("admin"), handler.BookTransactionHandler.CreateBookTransaction)
	// transactionRoute.PUT("/:id", middleware.RoleAuth("admin"), handler.BookTransactionHandler.UpdateBookTransaction)
	// transactionRoute.DELETE("/:id", middleware.RoleAuth("admin"), handler.BookTransactionHandler.DeleteBookTransaction)
	// transactionRoute.PATCH("/:id/status", middleware.RoleAuth("admin"), handler.BookTransactionHandler.UpdateStatus)
	// transactionRoute.POST("/:id/return", middleware.RoleAuth("admin"), handler.BookTransactionHandler.ReturnBook)

	return router
}
