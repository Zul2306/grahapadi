package config

import (
	"fmt"
	"log"
	"os"

	"inventory-backend/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds application configuration
type Config struct {
	Port      string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	JWTSecret string
	DB        *gorm.DB
}

// Load reads config from environment variables with defaults
func Load() *Config {
	// Load .env file (optional, won't fail if it doesn't exist)
	_ = godotenv.Load()

	return &Config{
		Port:      getEnv("PORT", "8080"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "postgres"),
		DBPass:    getEnv("DB_PASS", "password"),
		DBName:    getEnv("DB_NAME", "gudang"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
	}
}

// InitDB initializes database connection
func (c *Config) InitDB() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DBHost, c.DBUser, c.DBPass, c.DBName, c.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	c.DB = db
	log.Printf("✓ Database connected successfully to %s", c.DBName)
	return nil
}

// MigrateDB runs all database migrations
func (c *Config) MigrateDB() error {
	if c.DB == nil {
		return fmt.Errorf("database not initialized")
	}

	log.Printf("🔄 Checking and updating database schema...")

	// Run AutoMigrate to create/update tables with correct schema (preserves existing data)
	log.Printf("  - Creating tables with new schema...")
	if err := c.DB.AutoMigrate(
		&models.User{},
		&models.PasswordReset{},
		&models.Product{},
		&models.Gudang{},
		&models.Transaction{},
		&models.StockGudang{},
		&models.StockCard{},
		&models.StockOpname{},
	); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// Ensure produk table exists
	if !c.DB.Migrator().HasTable(&models.Produk{}) {
		log.Printf("  - Creating produk table...")
		c.DB.Migrator().CreateTable(&models.Produk{})
	}

	log.Printf("✓ Database migrations completed successfully")
	return nil
}

// TestConnection tests if database connection is working
func (c *Config) TestConnection() (bool, error) {
	if c.DB == nil {
		return false, fmt.Errorf("database not initialized")
	}

	sqlDB, err := c.DB.DB()
	if err != nil {
		return false, err
	}

	if err := sqlDB.Ping(); err != nil {
		return false, err
	}

	return true, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
