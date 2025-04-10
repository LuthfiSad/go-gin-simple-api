package main

import (
	"go-gin-simple-api/config"
	"go-gin-simple-api/lib"
	"go-gin-simple-api/router"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup DB
	db, err := config.SetupDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	cloudinary, err := lib.NewCloudinaryService(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to cloudinary: %v", err)
	}

	router := router.SetupRouter(db, cloudinary)

	// // Setup repositories
	// authRepo := repository.NewAuthRepository(db)
	// bookRepo := repository.NewBookRepository(db)
	// mediaRepo := repository.NewMediaRepository(db)
	// bookStockRepo := repository.NewBookStockRepository(db)
	// customerRepo := repository.NewCustomerRepository(db)
	// chargeRepo := repository.NewChargeRepository(db)
	// bookTransactionRepo := repository.NewBookTransactionRepository(db)

	// // Setup services
	// // cloudinaryService := lib.NewCloudinaryService(cfg)
	// authService := service.NewAuthService(authRepo)
	// bookService := service.NewBookService(bookRepo, mediaRepo)
	// mediaService := service.NewMediaService(mediaRepo, bookRepo, cloudinary)
	// bookStockService := service.NewBookStockService(bookStockRepo, bookRepo)
	// customerService := service.NewCustomerService(customerRepo, bookTransactionRepo)
	// chargeService := service.NewChargeService(chargeRepo, bookTransactionRepo, authRepo)
	// bookTransactionService := service.NewBookTransactionService(bookTransactionRepo, bookRepo, bookStockRepo, customerRepo)

	// // Setup handlers
	// authHandler := handler.NewAuthHandler(authService)
	// bookHandler := handler.NewBookHandler(bookService)
	// mediaHandler := handler.NewMediaHandler(mediaService)
	// bookStockHandler := handler.NewBookStockHandler(bookStockService)
	// customerHandler := handler.NewCustomerHandler(customerService)
	// chargeHandler := handler.NewChargeHandler(chargeService)
	// bookTransactionHandler := handler.NewBookTransactionHandler(bookTransactionService)

	// // Setup router
	// router := gin.Default()

	// api := router.Group("/api")
	// // Auth routes
	// api.POST("/register", authHandler.Register)
	// api.POST("/login", authHandler.Login)

	// // API routes with middleware
	// api.Use(middleware.JWTAuth(authRepo))

	// // Book routes
	// bookRoute := api.Group("/books")
	// bookRoute.GET("/", bookHandler.GetBooks)
	// bookRoute.GET("/:id", bookHandler.GetBookByID)
	// bookRoute.POST("/", middleware.RoleAuth("admin"), bookHandler.CreateBook)
	// bookRoute.PUT("/:id", middleware.RoleAuth("admin"), bookHandler.UpdateBook)
	// bookRoute.DELETE("/:id", middleware.RoleAuth("admin"), bookHandler.DeleteBook)
	// bookRoute.DELETE("/:id/cover", middleware.RoleAuth("admin"), bookHandler.DeleteBookCover)

	// // Media routes
	// media := api.Group("/media")
	// media.GET("/", middleware.RoleAuth("admin"), mediaHandler.GetMedias)
	// media.GET("/:id", middleware.RoleAuth("admin"), mediaHandler.GetMedia)
	// media.POST("/", middleware.RoleAuth("admin"), mediaHandler.UploadMedia)
	// media.DELETE("/:id", middleware.RoleAuth("admin"), mediaHandler.DeleteMedia)

	// // BookStock routes
	// bookStock := api.Group("/bookstocks")
	// bookStock.GET("", bookStockHandler.GetAll)
	// bookStock.GET("/:code", bookStockHandler.GetByCode)
	// bookStock.GET("/book/:book_id", bookStockHandler.GetByBookID)
	// bookStock.GET("/book/:book_id/available", bookStockHandler.GetAvailableByBookID)

	// // Protected routes (admin only)
	// bookStock.POST("", middleware.RoleAuth("admin"), bookStockHandler.Create)
	// bookStock.PUT("/:code", middleware.RoleAuth("admin"), bookStockHandler.Update)
	// bookStock.DELETE("/:code", middleware.RoleAuth("admin"), bookStockHandler.Delete)
	// bookStock.PATCH("/:code/status", middleware.RoleAuth("admin"), bookStockHandler.UpdateStatus)

	// // Customer routes
	// customerRoute := api.Group("/customers")
	// customerRoute.GET("", customerHandler.GetAll)
	// customerRoute.GET("/:id", customerHandler.GetByID)
	// customerRoute.GET("/:id/transactions", customerHandler.GetByIDWithTransactions)
	// customerRoute.GET("/code/:code", customerHandler.GetByCode)

	// // Protected customer routes (admin only)
	// customerRoute.POST("", middleware.RoleAuth("admin"), customerHandler.Create)
	// customerRoute.PUT("/:id", middleware.RoleAuth("admin"), customerHandler.Update)
	// customerRoute.DELETE("/:id", middleware.RoleAuth("admin"), customerHandler.Delete)

	// // Charge routes
	// chargeRoute := api.Group("/charges")
	// chargeRoute.GET("", chargeHandler.GetAll)
	// chargeRoute.GET("/:id", chargeHandler.GetByID)
	// chargeRoute.GET("/transaction/:transaction_id", chargeHandler.GetByBookTransactionID)
	// chargeRoute.GET("/user/:user_id", chargeHandler.GetByUserID)

	// // Protected charge routes
	// chargeRoute.POST("", chargeHandler.Create) // This already checks userData
	// chargeRoute.PUT("/:id", middleware.RoleAuth("admin"), chargeHandler.Update)
	// chargeRoute.DELETE("/:id", middleware.RoleAuth("admin"), chargeHandler.Delete)

	// // Book transaction routes
	// transactionRoute := api.Group("/transactions")
	// transactionRoute.GET("", bookTransactionHandler.GetAll)
	// transactionRoute.GET("/:id", bookTransactionHandler.GetByID)
	// transactionRoute.GET("/customer/:customer_id", bookTransactionHandler.GetByCustomerID)
	// transactionRoute.GET("/book/:book_id", bookTransactionHandler.GetByBookID)
	// transactionRoute.GET("/stock/:stock_code", bookTransactionHandler.GetByStockCode)
	// transactionRoute.GET("/overdue", bookTransactionHandler.GetOverdueTransactions)

	// // Protected book transaction routes (admin only)
	// transactionRoute.POST("", middleware.RoleAuth("admin"), bookTransactionHandler.Create)
	// transactionRoute.PUT("/:id", middleware.RoleAuth("admin"), bookTransactionHandler.Update)
	// transactionRoute.DELETE("/:id", middleware.RoleAuth("admin"), bookTransactionHandler.Delete)
	// transactionRoute.PATCH("/:id/status", middleware.RoleAuth("admin"), bookTransactionHandler.UpdateStatus)
	// transactionRoute.POST("/:id/return", middleware.RoleAuth("admin"), bookTransactionHandler.ReturnBook)

	// User routes
	// api.GET("/users", middleware.RoleAuth("admin"), userHandler.GetUsers)
	// api.GET("/users/:id", middleware.RoleAuth("admin"), userHandler.GetUserByID)
	// api.POST("/users", middleware.RoleAuth("admin"), userHandler.CreateUser)
	// api.PUT("/users/:id", middleware.RoleAuth("admin"), userHandler.UpdateUser)
	// api.DELETE("/users/:id", middleware.RoleAuth("admin"), userHandler.DeleteUser)

	// Start server
	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
