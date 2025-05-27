package main

import (
	"log"
	"os"

	"takadao/banking/internal/config"
	"takadao/banking/internal/handlers"
	"takadao/banking/internal/middleware"
	"takadao/banking/internal/repository"
	"takadao/banking/internal/routes"
	"takadao/banking/internal/service"

	_ "takadao/banking/docs" // This will be generated

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Banking API
// @version         1.0
// @description     A banking service API in Go using Gin framework.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      takadao-banking-api.onrender.com
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Set release mode in production
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := config.NewDatabaseConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	transactionService := service.NewTransactionService(transactionRepo)
	adminService := service.NewAdminService(userRepo, transactionRepo)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(os.Getenv("JWT_SECRET"))

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, authMiddleware)
	userHandler := handlers.NewUserHandler(userService, transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	adminHandler := handlers.NewAdminHandler(adminService)

	// Setup router
	router := routes.SetupRouter(authHandler, userHandler, transactionHandler, adminHandler, authMiddleware)

	// Add health check endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Add Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Configure trusted proxies
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
