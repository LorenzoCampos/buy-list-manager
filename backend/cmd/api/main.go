package main

import (
	"fmt"
	"log"

	"github.com/buylist-manager/backend/internal/config"
	"github.com/buylist-manager/backend/internal/database"
	"github.com/buylist-manager/backend/internal/handlers"
	"github.com/buylist-manager/backend/internal/repository"
	"github.com/buylist-manager/backend/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run database seeds (only in development)
	if cfg.Env == "development" {
		if err := database.Seed(db); err != nil {
			log.Printf("Warning: Failed to seed database: %v", err)
		} else {
			log.Println("Database seeded successfully")
		}
	}

	// Initialize repositories
	categoryRepo := repository.NewCategoryRepository(db)
	subcategoryRepo := repository.NewSubcategoryRepository(db)
	productRepo := repository.NewProductRepository(db)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "BuyList Manager API v1.0",
	})

	// Middleware
	app.Use(recover.New()) // Recover from panics
	app.Use(logger.New())  // Request logging

	// CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.FrontendURL,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH, OPTIONS",
	}))

	// Initialize services
	productService := services.NewProductService(productRepo, categoryRepo, subcategoryRepo)

	// Initialize handlers
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	subcategoryHandler := handlers.NewSubcategoryHandler(subcategoryRepo, categoryRepo)
	productHandler := handlers.NewProductHandler(productRepo, productService)

	// Routes
	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "BuyList Manager API is running",
		})
	})

	// Category routes
	categories := api.Group("/categories")
	categories.Get("/", categoryHandler.GetAll)
	categories.Get("/:id", categoryHandler.GetByID)
	categories.Post("/", categoryHandler.Create)
	categories.Put("/:id", categoryHandler.Update)
	categories.Delete("/:id", categoryHandler.Delete)

	// Subcategory routes
	subcategories := api.Group("/subcategories")
	subcategories.Get("/", subcategoryHandler.GetAll)           // GET /api/v1/subcategories?category_id=1
	subcategories.Get("/:id", subcategoryHandler.GetByID)       // GET /api/v1/subcategories/1
	subcategories.Post("/", subcategoryHandler.Create)          // POST /api/v1/subcategories
	subcategories.Put("/:id", subcategoryHandler.Update)        // PUT /api/v1/subcategories/1
	subcategories.Delete("/:id", subcategoryHandler.Delete)     // DELETE /api/v1/subcategories/1

	// Product routes
	products := api.Group("/products")
	products.Get("/", productHandler.GetAll)                    // GET /api/v1/products?pending=true&category_id=1
	products.Get("/stats", productHandler.GetStats)             // GET /api/v1/products/stats
	products.Get("/:id", productHandler.GetByID)                // GET /api/v1/products/1
	products.Post("/", productHandler.Create)                   // POST /api/v1/products
	products.Put("/:id", productHandler.Update)                 // PUT /api/v1/products/1
	products.Delete("/:id", productHandler.Delete)              // DELETE /api/v1/products/1

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🚀 Server starting on http://localhost%s", addr)
	log.Printf("📖 Environment: %s", cfg.Env)
	log.Printf("🗄️  Database: %s", cfg.DBDriver)

	if err := app.Listen(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
