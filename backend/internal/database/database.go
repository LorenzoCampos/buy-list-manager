package database

import (
	"fmt"

	"github.com/buylist-manager/backend/internal/config"
	"github.com/buylist-manager/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establishes a database connection and runs migrations
func Connect(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	// Choose database driver (only PostgreSQL supported for now)
	if cfg.DBDriver != "postgres" {
		return nil, fmt.Errorf("only postgres driver is supported, got: %s", cfg.DBDriver)
	}
	dialector = postgres.Open(cfg.GetDatabaseDSN())

	// Configure GORM
	gormConfig := &gorm.Config{}
	if cfg.Env == "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	// Open connection
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// migrate runs all database migrations
func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Category{},
		&models.Subcategory{},
		&models.Product{},
	)
}

// Seed populates the database with initial data
func Seed(db *gorm.DB) error {
	// Check if categories already exist
	var count int64
	db.Model(&models.Category{}).Count(&count)
	if count > 0 {
		return nil // Already seeded
	}

	// Create categories
	categories := []models.Category{
		{Name: "Compra Única", Type: "one_time"},
		{Name: "Suscripción Mensual", Type: "recurring"},
		{Name: "Suscripción Anual", Type: "recurring"},
	}

	for i := range categories {
		if err := db.Create(&categories[i]).Error; err != nil {
			return fmt.Errorf("failed to seed categories: %w", err)
		}
	}

	// Create subcategories
	subcategories := []models.Subcategory{
		{CategoryID: 1, Name: "Reparación de Electrónicos"},
		{CategoryID: 1, Name: "Trabajo/Productividad"},
		{CategoryID: 1, Name: "Gaming"},
		{CategoryID: 1, Name: "Hogar"},
		{CategoryID: 2, Name: "IA y Herramientas"},
		{CategoryID: 2, Name: "Entretenimiento"},
		{CategoryID: 3, Name: "Software Profesional"},
	}

	for i := range subcategories {
		if err := db.Create(&subcategories[i]).Error; err != nil {
			return fmt.Errorf("failed to seed subcategories: %w", err)
		}
	}

	return nil
}
