package models

import (
	"time"

	"gorm.io/gorm"
)

// Product represents an item to buy or a subscription
type Product struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"size:255;not null" json:"name"`
	Description    string         `gorm:"type:text" json:"description"`
	BasePrice      float64        `gorm:"type:decimal(10,2);not null" json:"base_price"`
	ShippingCost   float64        `gorm:"type:decimal(10,2);default:0" json:"shipping_cost"`
	Taxes          float64        `gorm:"type:decimal(10,2);default:0" json:"taxes"`
	TotalPrice     float64        `gorm:"type:decimal(10,2)" json:"total_price"` // Calculated field
	SourceURL      string         `gorm:"size:500" json:"source_url"`
	PriceDate      *time.Time     `json:"price_date"`
	CategoryID     uint           `gorm:"not null" json:"category_id"`
	SubcategoryID  uint           `gorm:"not null" json:"subcategory_id"`
	RecurrenceType *string        `gorm:"size:20" json:"recurrence_type"` // null, "monthly", "yearly"
	IsPurchased    bool           `gorm:"default:false" json:"is_purchased"`
	PurchaseDate   *time.Time     `json:"purchase_date"`
	Notes          string         `gorm:"type:text" json:"notes"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relationships
	Category    *Category    `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Subcategory *Subcategory `gorm:"foreignKey:SubcategoryID" json:"subcategory,omitempty"`
}

// TableName specifies the table name for GORM
func (Product) TableName() string {
	return "products"
}

// BeforeSave is a GORM hook that calculates TotalPrice before saving
func (p *Product) BeforeSave(tx *gorm.DB) error {
	p.TotalPrice = p.BasePrice + p.ShippingCost + p.Taxes
	return nil
}

// IsValidRecurrenceType checks if the recurrence type is valid
func (p *Product) IsValidRecurrenceType() bool {
	if p.RecurrenceType == nil {
		return true // null is valid
	}
	return *p.RecurrenceType == "monthly" || *p.RecurrenceType == "yearly"
}
