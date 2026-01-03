package services

import (
	"errors"

	"github.com/buylist-manager/backend/internal/models"
	"github.com/buylist-manager/backend/internal/repository"
)

// ProductService handles business logic for products
type ProductService interface {
	CreateProduct(product *models.Product) error
	GetTotalPendingCost() (float64, error)
	GetMonthlyRecurringCost() (float64, error)
	GetYearlyRecurringCost() (float64, error)
}

// productService is the concrete implementation
type productService struct {
	productRepo      repository.ProductRepository
	categoryRepo     repository.CategoryRepository
	subcategoryRepo  repository.SubcategoryRepository
}

// NewProductService creates a new instance of ProductService
func NewProductService(
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	subcategoryRepo repository.SubcategoryRepository,
) ProductService {
	return &productService{
		productRepo:     productRepo,
		categoryRepo:    categoryRepo,
		subcategoryRepo: subcategoryRepo,
	}
}

// CreateProduct creates a product with validations
func (s *productService) CreateProduct(product *models.Product) error {
	// Validar que la categoría existe
	category, err := s.categoryRepo.FindByID(product.CategoryID)
	if err != nil {
		return errors.New("category not found")
	}

	// Validar que la subcategoría existe
	_, err = s.subcategoryRepo.FindByID(product.SubcategoryID)
	if err != nil {
		return errors.New("subcategory not found")
	}

	// Validar coherencia entre category.type y recurrence_type
	// Si la categoría es "one_time", recurrence_type debe ser null
	// Si la categoría es "recurring", recurrence_type debe ser "monthly" o "yearly"
	if category.Type == "one_time" {
		if product.RecurrenceType != nil {
			return errors.New("one-time purchases cannot have recurrence type")
		}
	} else if category.Type == "recurring" {
		if product.RecurrenceType == nil {
			return errors.New("recurring purchases must have recurrence type")
		}
		if *product.RecurrenceType != "monthly" && *product.RecurrenceType != "yearly" {
			return errors.New("recurrence type must be 'monthly' or 'yearly'")
		}
	}

	// El cálculo de total_price se hace automáticamente en el hook BeforeSave del modelo
	return s.productRepo.Create(product)
}

// GetTotalPendingCost calcula el total de productos pendientes de compra (one-time)
// Laravel: Product::where('is_purchased', false)->whereHas('category', fn($q) => $q->where('type', 'one_time'))->sum('total_price')
func (s *productService) GetTotalPendingCost() (float64, error) {
	products, err := s.productRepo.FindPending()
	if err != nil {
		return 0, err
	}

	var total float64
	for _, product := range products {
		// Solo sumar productos de compra única
		if product.Category != nil && product.Category.Type == "one_time" {
			total += product.TotalPrice
		}
	}

	return total, nil
}

// GetMonthlyRecurringCost calcula el gasto mensual en suscripciones
func (s *productService) GetMonthlyRecurringCost() (float64, error) {
	allProducts, err := s.productRepo.FindAll()
	if err != nil {
		return 0, err
	}

	var total float64
	for _, product := range allProducts {
		// Solo productos recurring y no comprados (o purchased si son suscripciones activas)
		if product.Category != nil && product.Category.Type == "recurring" && product.RecurrenceType != nil {
			if *product.RecurrenceType == "monthly" {
				total += product.TotalPrice
			} else if *product.RecurrenceType == "yearly" {
				// Dividir el costo anual en 12 meses
				total += product.TotalPrice / 12
			}
		}
	}

	return total, nil
}

// GetYearlyRecurringCost calcula el gasto anual en suscripciones
func (s *productService) GetYearlyRecurringCost() (float64, error) {
	allProducts, err := s.productRepo.FindAll()
	if err != nil {
		return 0, err
	}

	var total float64
	for _, product := range allProducts {
		if product.Category != nil && product.Category.Type == "recurring" && product.RecurrenceType != nil {
			if *product.RecurrenceType == "yearly" {
				total += product.TotalPrice
			} else if *product.RecurrenceType == "monthly" {
				// Multiplicar el costo mensual por 12
				total += product.TotalPrice * 12
			}
		}
	}

	return total, nil
}
