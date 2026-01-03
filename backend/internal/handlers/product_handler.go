package handlers

import (
	"strconv"
	"time"

	"github.com/buylist-manager/backend/internal/models"
	"github.com/buylist-manager/backend/internal/repository"
	"github.com/buylist-manager/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

// ProductHandler handles HTTP requests for products
type ProductHandler struct {
	repo    repository.ProductRepository
	service services.ProductService
}

// NewProductHandler creates a new ProductHandler
func NewProductHandler(repo repository.ProductRepository, service services.ProductService) *ProductHandler {
	return &ProductHandler{
		repo:    repo,
		service: service,
	}
}

// GetAll retrieves all products with optional filters
func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	// Filtros opcionales por query params
	categoryIDStr := c.Query("category_id")
	subcategoryIDStr := c.Query("subcategory_id")
	pendingStr := c.Query("pending")

	// Filtrar por pending (productos no comprados)
	if pendingStr == "true" {
		products, err := h.repo.FindPending()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch products",
			})
		}
		return c.JSON(products)
	}

	// Filtrar por categoría
	if categoryIDStr != "" {
		categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid category_id parameter",
			})
		}

		products, err := h.repo.FindByCategoryID(uint(categoryID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch products",
			})
		}
		return c.JSON(products)
	}

	// Filtrar por subcategoría
	if subcategoryIDStr != "" {
		subcategoryID, err := strconv.ParseUint(subcategoryIDStr, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid subcategory_id parameter",
			})
		}

		products, err := h.repo.FindBySubcategoryID(uint(subcategoryID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch products",
			})
		}
		return c.JSON(products)
	}

	// Sin filtros, traer todos
	products, err := h.repo.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch products",
		})
	}

	return c.JSON(products)
}

// GetByID retrieves a single product by ID
func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	product, err := h.repo.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.JSON(product)
}

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	Name           string  `json:"name" validate:"required,min=1,max=255"`
	Description    string  `json:"description"`
	BasePrice      float64 `json:"base_price" validate:"required,min=0"`
	ShippingCost   float64 `json:"shipping_cost" validate:"min=0"`
	Taxes          float64 `json:"taxes" validate:"min=0"`
	SourceURL      string  `json:"source_url"`
	CategoryID     uint    `json:"category_id" validate:"required"`
	SubcategoryID  uint    `json:"subcategory_id" validate:"required"`
	RecurrenceType *string `json:"recurrence_type"` // "monthly" o "yearly" o null
	Notes          string  `json:"notes"`
}

// Create creates a new product
func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var req CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	now := time.Now()
	product := &models.Product{
		Name:           req.Name,
		Description:    req.Description,
		BasePrice:      req.BasePrice,
		ShippingCost:   req.ShippingCost,
		Taxes:          req.Taxes,
		SourceURL:      req.SourceURL,
		PriceDate:      &now,
		CategoryID:     req.CategoryID,
		SubcategoryID:  req.SubcategoryID,
		RecurrenceType: req.RecurrenceType,
		Notes:          req.Notes,
		IsPurchased:    false,
	}

	// Usar el Service que tiene las validaciones de negocio
	if err := h.service.CreateProduct(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

// UpdateProductRequest represents the request body for updating a product
type UpdateProductRequest struct {
	Name           string  `json:"name" validate:"required,min=1,max=255"`
	Description    string  `json:"description"`
	BasePrice      float64 `json:"base_price" validate:"required,min=0"`
	ShippingCost   float64 `json:"shipping_cost" validate:"min=0"`
	Taxes          float64 `json:"taxes" validate:"min=0"`
	SourceURL      string  `json:"source_url"`
	CategoryID     uint    `json:"category_id" validate:"required"`
	SubcategoryID  uint    `json:"subcategory_id" validate:"required"`
	RecurrenceType *string `json:"recurrence_type"`
	IsPurchased    bool    `json:"is_purchased"`
	Notes          string  `json:"notes"`
}

// Update updates an existing product
func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	// Check if product exists
	product, err := h.repo.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	var req UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	product.Name = req.Name
	product.Description = req.Description
	product.BasePrice = req.BasePrice
	product.ShippingCost = req.ShippingCost
	product.Taxes = req.Taxes
	product.SourceURL = req.SourceURL
	product.CategoryID = req.CategoryID
	product.SubcategoryID = req.SubcategoryID
	product.RecurrenceType = req.RecurrenceType
	product.IsPurchased = req.IsPurchased
	product.Notes = req.Notes

	// Si se marca como comprado, guardar la fecha
	if req.IsPurchased && product.PurchaseDate == nil {
		now := time.Now()
		product.PurchaseDate = &now
	}

	if err := h.repo.Update(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update product",
		})
	}

	return c.JSON(product)
}

// Delete deletes a product by ID
func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetStats returns statistics about products (totals, monthly cost, etc.)
func (h *ProductHandler) GetStats(c *fiber.Ctx) error {
	totalPending, err := h.service.GetTotalPendingCost()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate statistics",
		})
	}

	monthlyCost, err := h.service.GetMonthlyRecurringCost()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate statistics",
		})
	}

	yearlyCost, err := h.service.GetYearlyRecurringCost()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate statistics",
		})
	}

	return c.JSON(fiber.Map{
		"total_pending_one_time": totalPending,
		"monthly_recurring_cost": monthlyCost,
		"yearly_recurring_cost":  yearlyCost,
	})
}
