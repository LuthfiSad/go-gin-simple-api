package main

import (
	"go-gin-simple-api/config"
	"go-gin-simple-api/handler"
	"go-gin-simple-api/middleware"
	"go-gin-simple-api/repository"
	"go-gin-simple-api/service"
	"log"

	"github.com/gin-gonic/gin"
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

	// Setup repositories
	authRepo := repository.NewAuthRepository(db)
	bookRepo := repository.NewBookRepository(db)

	// Setup services
	// cloudinaryService := lib.NewCloudinaryService(cfg)
	authService := service.NewAuthService(authRepo)
	bookService := service.NewBookService(bookRepo)

	// Setup handlers
	authHandler := handler.NewAuthHandler(authService)
	bookHandler := handler.NewBookHandler(bookService)

	// Setup router
	router := gin.Default()

	api := router.Group("/api")
	// Auth routes
	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)

	// API routes with middleware
	api.Use(middleware.JWTAuth(authRepo))

	// Book routes
	api.GET("/books", bookHandler.GetBooks)
	api.GET("/books/:id", bookHandler.GetBookByID)
	api.POST("/books", middleware.RoleAuth("admin"), bookHandler.CreateBook)
	api.PUT("/books/:id", middleware.RoleAuth("admin"), bookHandler.UpdateBook)
	api.DELETE("/books/:id", middleware.RoleAuth("admin"), bookHandler.DeleteBook)
	api.DELETE("/books/:id/cover", middleware.RoleAuth("admin"), bookHandler.DeleteBookCover)

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
