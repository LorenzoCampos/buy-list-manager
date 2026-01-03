package repository

import (
	"errors"

	"github.com/buylist-manager/backend/internal/models"
	"gorm.io/gorm"
)

// CategoryRepository defines the interface for category data operations
type CategoryRepository interface {
	Create(category *models.Category) error
	FindByID(id uint) (*models.Category, error)
	FindAll() ([]*models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
}

// categoryRepository is the concrete implementation
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new instance of CategoryRepository
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// Create inserts a new category into the database
func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// FindByID retrieves a category by its ID
func (r *categoryRepository) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// FindAll retrieves all categories
func (r *categoryRepository) FindAll() ([]*models.Category, error) {
	var categories []*models.Category
	err := r.db.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// Update updates an existing category
func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete deletes a category by ID (soft delete if using GORM soft delete)
func (r *categoryRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Category{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}
