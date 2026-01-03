package repository

import (
	"errors"

	"github.com/buylist-manager/backend/internal/models"
	"gorm.io/gorm"
)

// SubcategoryRepository defines the interface for subcategory data operations
type SubcategoryRepository interface {
	Create(subcategory *models.Subcategory) error
	FindByID(id uint) (*models.Subcategory, error)
	FindAll() ([]*models.Subcategory, error)
	FindByCategoryID(categoryID uint) ([]*models.Subcategory, error)
	Update(subcategory *models.Subcategory) error
	Delete(id uint) error
}

// subcategoryRepository is the concrete implementation
type subcategoryRepository struct {
	db *gorm.DB
}

// NewSubcategoryRepository creates a new instance of SubcategoryRepository
func NewSubcategoryRepository(db *gorm.DB) SubcategoryRepository {
	return &subcategoryRepository{db: db}
}

// Create inserts a new subcategory into the database
func (r *subcategoryRepository) Create(subcategory *models.Subcategory) error {
	return r.db.Create(subcategory).Error
}

// FindByID retrieves a subcategory by its ID
func (r *subcategoryRepository) FindByID(id uint) (*models.Subcategory, error) {
	var subcategory models.Subcategory
	// Preload carga la relación Category (como ->with('category') en Laravel)
	err := r.db.Preload("Category").First(&subcategory, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subcategory not found")
		}
		return nil, err
	}
	return &subcategory, nil
}

// FindAll retrieves all subcategories
func (r *subcategoryRepository) FindAll() ([]*models.Subcategory, error) {
	var subcategories []*models.Subcategory
	// Preload Category para cada subcategory
	err := r.db.Preload("Category").Find(&subcategories).Error
	if err != nil {
		return nil, err
	}
	return subcategories, nil
}

// FindByCategoryID retrieves all subcategories for a specific category
// Equivalente en Laravel: Subcategory::where('category_id', $categoryId)->get()
func (r *subcategoryRepository) FindByCategoryID(categoryID uint) ([]*models.Subcategory, error) {
	var subcategories []*models.Subcategory
	err := r.db.Where("category_id = ?", categoryID).Find(&subcategories).Error
	if err != nil {
		return nil, err
	}
	return subcategories, nil
}

// Update updates an existing subcategory
func (r *subcategoryRepository) Update(subcategory *models.Subcategory) error {
	return r.db.Save(subcategory).Error
}

// Delete deletes a subcategory by ID
func (r *subcategoryRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Subcategory{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("subcategory not found")
	}
	return nil
}
