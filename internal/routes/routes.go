package routes

import (
	"takadao/banking/internal/handlers"
	"takadao/banking/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all the routes for the application
func SetupRouter(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	transactionHandler *handlers.TransactionHandler,
	adminHandler *handlers.AdminHandler,
	authMiddleware *middleware.AuthMiddleware,
) *gin.Engine {
	router := gin.Default()

	// API version 1
	api := router.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			// User auth routes
			userAuth := auth.Group("/user")
			{
				userAuth.POST("/register", authHandler.RegisterUser)
				userAuth.POST("/login", authHandler.UserLogin)
			}

			// Admin auth routes
			adminAuth := auth.Group("/admin")
			{
				adminAuth.POST("/login", authHandler.AdminLogin)
				adminAuth.POST("/register", authMiddleware.RequireAuth(), authHandler.RegisterAdmin)
			}
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(authMiddleware.RequireAuth())
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.GetMe)
				users.PUT("/me", userHandler.UpdateMe)
				users.GET("/balance", userHandler.GetBalances)
			}

			// Transaction routes
			transactions := protected.Group("/transactions")
			{
				transactions.GET("/me", transactionHandler.ListMyTransactions)
				transactions.GET("/me/:id", transactionHandler.GetMyTransaction)
				transactions.POST("/deposit", transactionHandler.Deposit)
				transactions.POST("/withdraw", transactionHandler.Withdraw)
				transactions.POST("/transfer", transactionHandler.Transfer)
			}

			// Admin routes
			admin := protected.Group("/admin")
			admin.Use(authMiddleware.RequireAdmin())
			{
				// User management
				admin.GET("/users", userHandler.ListUsers)
				admin.GET("/users/:id", userHandler.GetUser)
				admin.PUT("/users/:id", userHandler.UpdateUser)
				admin.DELETE("/users/:id", userHandler.DeleteUser)

				// Transaction management
				admin.GET("/transactions", adminHandler.ListAllTransactions)
				admin.GET("/transactions/:id", transactionHandler.GetTransaction)
			}
		}
	}

	return router
}
