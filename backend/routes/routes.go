package routes

import (
	"inventory-backend/config"
	"inventory-backend/controllers"
	"inventory-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes all routes
func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Apply global middleware
	r.Use(middleware.CORS())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Inventory API is running"})
	})

	// API v1 group
	api := r.Group("/api/v1")
	{
		// Database test endpoint (public)
		api.GET("/db/test", func(c *gin.Context) {
			connected, err := cfg.TestConnection()
			if err != nil {
				c.JSON(500, gin.H{
					"status": "error",
					"message": "Database connection failed",
					"error": err.Error(),
				})
				return
			}
			
			if connected {
				c.JSON(200, gin.H{
					"status": "connected",
					"message": "Database connection is working",
					"database": "gudang",
					"host": cfg.DBHost,
					"port": cfg.DBPort,
				})
			} else {
				c.JSON(500, gin.H{
					"status": "error",
					"message": "Database connection test failed",
				})
			}
		})

		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", func(c *gin.Context) {
				c.Set("config", cfg)
				controllers.Login(c)
			})
			auth.POST("/logout", controllers.Logout)
			auth.POST("/register", func(c *gin.Context) {
				c.Set("config", cfg)
				controllers.Register(c)
			})
			auth.POST("/forgot-password", func(c *gin.Context) {
				c.Set("config", cfg)
				controllers.ForgotPassword(c)
			})
			auth.POST("/reset-password", func(c *gin.Context) {
				c.Set("config", cfg)
				controllers.ResetPassword(c)
			})
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			// User management
			users := protected.Group("/users")
			{
				users.GET("", controllers.GetUsers)
				users.GET("/:id", controllers.GetUser)
				users.POST("", controllers.CreateUser)
				users.DELETE("/:id", controllers.DeleteUser)
			}

			// Product / Item management
			products := protected.Group("/products")
			{
				products.GET("", controllers.GetProducts)
				products.GET("/:id", controllers.GetProduct)
				products.POST("", controllers.CreateProduct)
				products.PUT("/:id", controllers.UpdateProduct)
				products.DELETE("/:id", controllers.DeleteProduct)
			}

			// Stock / Opname
			stock := protected.Group("/stock")
			{
				stock.GET("", controllers.GetStockCards)
				stock.GET("/:id", controllers.GetStockCard)
				stock.POST("/opname", controllers.CreateOpname)
			}
		}
	}

	return r
}
