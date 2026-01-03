package handlers

import (
	"strconv"

	"github.com/buylist-manager/backend/internal/models"
	"github.com/buylist-manager/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
)

// CategoryHandler handles HTTP requests for categories
type CategoryHandler struct {
	repo repository.CategoryRepository
}

// NewCategoryHandler creates a new CategoryHandler
func NewCategoryHandler(repo repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

// GetAll retrieves all categories
// @Summary Get all categories
// @Description Get a list of all categories
// @Tags categories
// @Produce json
// @Success 200 {array} models.Category
// @Router /api/v1/categories [get]
func (h *CategoryHandler) GetAll(c *fiber.Ctx) error {
	categories, err := h.repo.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch categories",
		})
	}

	return c.JSON(categories)
}

// GetByID retrieves a single category by ID
// @Summary Get category by ID
// @Description Get a single category by its ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.Category
// @Router /api/v1/categories/{id} [get]
func (h *CategoryHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	category, err := h.repo.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	return c.JSON(category)
}

// CreateCategoryRequest represents the request body for creating a category
type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	Type string `json:"type" validate:"required,oneof=one_time recurring"`
}

// Create creates a new category
// @Summary Create a new category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body CreateCategoryRequest true "Category data"
// @Success 201 {object} models.Category
// @Router /api/v1/categories [post]
func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	var req CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate category type
	if req.Type != "one_time" && req.Type != "recurring" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category type. Must be 'one_time' or 'recurring'",
		})
	}

	category := &models.Category{
		Name: req.Name,
		Type: req.Type,
	}

	if err := h.repo.Create(category); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create category",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

// UpdateCategoryRequest represents the request body for updating a category
type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	Type string `json:"type" validate:"required,oneof=one_time recurring"`
}

// Update updates an existing category
// @Summary Update a category
// @Description Update a category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body UpdateCategoryRequest true "Category data"
// @Success 200 {object} models.Category
// @Router /api/v1/categories/{id} [put]
func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	// Check if category exists
	category, err := h.repo.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	var req UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate category type
	if req.Type != "one_time" && req.Type != "recurring" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category type. Must be 'one_time' or 'recurring'",
		})
	}

	// Update fields
	category.Name = req.Name
	category.Type = req.Type

	if err := h.repo.Update(category); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update category",
		})
	}

	return c.JSON(category)
}

// Delete deletes a category by ID
// @Summary Delete a category
// @Description Delete a category by ID
// @Tags categories
// @Param id path int true "Category ID"
// @Success 204
// @Router /api/v1/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
