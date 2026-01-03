package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a main category (one-time purchase or recurring subscription)
type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Type      string         `gorm:"size:20;not null" json:"type"` // "one_time" or "recurring"
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relationships
	Subcategories []Subcategory `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"subcategories,omitempty"`
}

// TableName specifies the table name for GORM
func (Category) TableName() string {
	return "categories"
}

// IsValidType checks if the category type is valid
func (c *Category) IsValidType() bool {
	return c.Type == "one_time" || c.Type == "recurring"
}
