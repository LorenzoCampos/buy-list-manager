package handlers

import (
	"strconv"

	"github.com/buylist-manager/backend/internal/models"
	"github.com/buylist-manager/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
)

// SubcategoryHandler handles HTTP requests for subcategories
type SubcategoryHandler struct {
	repo         repository.SubcategoryRepository
	categoryRepo repository.CategoryRepository // Para validar que la categoría existe
}

// NewSubcategoryHandler creates a new SubcategoryHandler
func NewSubcategoryHandler(repo repository.SubcategoryRepository, categoryRepo repository.CategoryRepository) *SubcategoryHandler {
	return &SubcategoryHandler{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

// GetAll retrieves all subcategories
func (h *SubcategoryHandler) GetAll(c *fiber.Ctx) error {
	// Si viene el query param ?category_id=X, filtramos
	categoryIDStr := c.Query("category_id")
	
	if categoryIDStr != "" {
		categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid category_id parameter",
			})
		}
		
		subcategories, err := h.repo.FindByCategoryID(uint(categoryID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch subcategories",
			})
		}
		return c.JSON(subcategories)
	}
	
	// Sin filtro, traer todas
	subcategories, err := h.repo.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch subcategories",
		})
	}

	return c.JSON(subcategories)
}

// GetByID retrieves a single subcategory by ID
func (h *SubcategoryHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid subcategory ID",
		})
	}

	subcategory, err := h.repo.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Subcategory not found",
		})
	}

	return c.JSON(subcategory)
}

// CreateSubcategoryRequest represents the request body for creating a subcategory
type CreateSubcategoryRequest struct {
	CategoryID uint   `json:"category_id" validate:"required"`
	Name       string `json:"name" validate:"required,min=1,max=100"`
}

// Create creates a new subcategory
func (h *SubcategoryHandler) Create(c *fiber.Ctx) error {
	var req CreateSubcategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validar que la categoría existe (como en Laravel con exists:categories,id)
	_, err := h.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	subcategory := &models.Subcategory{
		CategoryID: req.CategoryID,
		Name:       req.Name,
	}

	if err := h.repo.Create(subcategory); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create subcategory",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(subcategory)
}

// UpdateSubcategoryRequest represents the request body for updating a subcategory
type UpdateSubcategoryRequest struct {
	CategoryID uint   `json:"category_id" validate:"required"`
	Name       string `json:"name" validate:"required,min=1,max=100"`
}

// Update updates an existing subcategory
func (h *SubcategoryHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid subcategory ID",
		})
	}

	// Check if subcategory exists
	subcategory, err := h.repo.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Subcategory not found",
		})
	}

	var req UpdateSubcategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validar que la nueva categoría existe
	_, err = h.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	// Update fields
	subcategory.CategoryID = req.CategoryID
	subcategory.Name = req.Name

	if err := h.repo.Update(subcategory); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update subcategory",
		})
	}

	return c.JSON(subcategory)
}

// Delete deletes a subcategory by ID
func (h *SubcategoryHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid subcategory ID",
		})
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Subcategory not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
