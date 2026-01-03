package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBDriver   string // "postgres" or "sqlite"
	DBDSN      string // For SQLite

	// Server
	Port string
	Env  string

	// CORS
	FrontendURL string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not found in production)
	godotenv.Load()

	cfg := &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "buylist_db"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		DBDriver:    getEnv("DB_DRIVER", "postgres"),
		DBDSN:       getEnv("DB_DSN", "./buylist.db"),
		Port:        getEnv("PORT", "8080"),
		Env:         getEnv("ENV", "development"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
	}

	return cfg, nil
}

// GetDatabaseDSN returns the database connection string
func (c *Config) GetDatabaseDSN() string {
	if c.DBDriver == "sqlite" {
		return c.DBDSN
	}

	// PostgreSQL DSN (omit password if empty for peer authentication)
	if c.DBPassword == "" {
		return fmt.Sprintf(
			"host=%s user=%s dbname=%s port=%s sslmode=%s",
			c.DBHost, c.DBUser, c.DBName, c.DBPort, c.DBSSLMode,
		)
	}
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.DBSSLMode,
	)
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
