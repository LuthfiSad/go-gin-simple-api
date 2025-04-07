package main

import (
	"go-gin-simple-api/config"
	"go-gin-simple-api/handler"
	"go-gin-simple-api/lib"
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

	cloudinary, err := lib.NewCloudinaryService(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to cloudinary: %v", err)
	}

	// Setup repositories
	authRepo := repository.NewAuthRepository(db)
	bookRepo := repository.NewBookRepository(db)
	mediaRepo := repository.NewMediaRepository(db)

	// Setup services
	// cloudinaryService := lib.NewCloudinaryService(cfg)
	authService := service.NewAuthService(authRepo)
	bookService := service.NewBookService(bookRepo, mediaRepo)
	mediaService := service.NewMediaService(mediaRepo, bookRepo, cloudinary)

	// Setup handlers
	authHandler := handler.NewAuthHandler(authService)
	bookHandler := handler.NewBookHandler(bookService)
	mediaHandler := handler.NewMediaHandler(mediaService)

	// Setup router
	router := gin.Default()

	api := router.Group("/api")
	// Auth routes
	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)

	// API routes with middleware
	api.Use(middleware.JWTAuth(authRepo))

	// Book routes
	bookRoute := api.Group("/books")
	bookRoute.GET("/", bookHandler.GetBooks)
	bookRoute.GET("/:id", bookHandler.GetBookByID)
	bookRoute.POST("/", middleware.RoleAuth("admin"), bookHandler.CreateBook)
	bookRoute.PUT("/:id", middleware.RoleAuth("admin"), bookHandler.UpdateBook)
	bookRoute.DELETE("/:id", middleware.RoleAuth("admin"), bookHandler.DeleteBook)
	bookRoute.DELETE("/:id/cover", middleware.RoleAuth("admin"), bookHandler.DeleteBookCover)

	// Media routes
	media := api.Group("/media")
	media.GET("/", middleware.RoleAuth("admin"), mediaHandler.GetMedias)
	media.GET("/:id", middleware.RoleAuth("admin"), mediaHandler.GetMedia)
	media.POST("/", middleware.RoleAuth("admin"), mediaHandler.UploadMedia)
	media.DELETE("/:id", middleware.RoleAuth("admin"), mediaHandler.DeleteMedia)

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
