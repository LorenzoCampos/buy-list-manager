package models

import (
	"time"

	"gorm.io/gorm"
)

// Subcategory represents a sub-category within a main category
type Subcategory struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CategoryID uint           `gorm:"not null" json:"category_id"`
	Name       string         `gorm:"size:100;not null" json:"name"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relationships
	Category *Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"category,omitempty"`
	Products []Product `gorm:"foreignKey:SubcategoryID" json:"products,omitempty"`
}

// TableName specifies the table name for GORM
func (Subcategory) TableName() string {
	return "subcategories"
}
