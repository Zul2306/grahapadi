package main

import (
	"inventory-backend/config"
	"inventory-backend/routes"
	"log"
)

var globalConfig *config.Config

func main() {
	// Load configuration
	cfg := config.Load()
	globalConfig = cfg

	// Initialize database connection
	if err := cfg.InitDB(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// Run migrations
	if err := cfg.MigrateDB(); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	// Test connection
	if connected, err := cfg.TestConnection(); err != nil {
		log.Fatalf("Database connection test failed: %v", err)
	} else if connected {
		log.Printf("✓ Database connection verified")
	}

	// Setup router
	r := routes.SetupRouter(cfg)

	log.Printf("🚀 Server running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func GetConfig() *config.Config {
	return globalConfig
}
