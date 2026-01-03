package repository

import (
	"errors"

	"github.com/buylist-manager/backend/internal/models"
	"gorm.io/gorm"
)

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	Create(product *models.Product) error
	FindByID(id uint) (*models.Product, error)
	FindAll() ([]*models.Product, error)
	FindByCategoryID(categoryID uint) ([]*models.Product, error)
	FindBySubcategoryID(subcategoryID uint) ([]*models.Product, error)
	FindPending() ([]*models.Product, error) // Productos no comprados
	Update(product *models.Product) error
	Delete(id uint) error
}

// productRepository is the concrete implementation
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new instance of ProductRepository
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create inserts a new product into the database
func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

// FindByID retrieves a product by its ID
func (r *productRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	// Preload Category y Subcategory (como ->with(['category', 'subcategory']) en Laravel)
	err := r.db.Preload("Category").Preload("Subcategory").First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

// FindAll retrieves all products
func (r *productRepository) FindAll() ([]*models.Product, error) {
	var products []*models.Product
	err := r.db.Preload("Category").Preload("Subcategory").
		Order("created_at DESC").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// FindByCategoryID retrieves all products for a specific category
// Laravel: Product::where('category_id', $categoryId)->get()
func (r *productRepository) FindByCategoryID(categoryID uint) ([]*models.Product, error) {
	var products []*models.Product
	err := r.db.Preload("Category").Preload("Subcategory").
		Where("category_id = ?", categoryID).
		Order("created_at DESC").
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// FindBySubcategoryID retrieves all products for a specific subcategory
func (r *productRepository) FindBySubcategoryID(subcategoryID uint) ([]*models.Product, error) {
	var products []*models.Product
	err := r.db.Preload("Category").Preload("Subcategory").
		Where("subcategory_id = ?", subcategoryID).
		Order("created_at DESC").
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// FindPending retrieves all products that haven't been purchased yet
// Laravel: Product::where('is_purchased', false)->get()
func (r *productRepository) FindPending() ([]*models.Product, error) {
	var products []*models.Product
	err := r.db.Preload("Category").Preload("Subcategory").
		Where("is_purchased = ?", false).
		Order("created_at DESC").
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// Update updates an existing product
func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

// Delete deletes a product by ID
func (r *productRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}
