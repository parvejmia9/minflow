package models

import (
	"time"

	"gorm.io/gorm"
)

// Expense represents an expense entry in the system
type Expense struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	CategoryID  uint           `gorm:"not null" json:"category_id"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Unit        float64        `gorm:"not null" json:"unit"`
	PerUnitCost float64        `gorm:"not null;type:decimal(10,2)" json:"per_unit_cost"`
	Total       float64        `gorm:"not null;type:decimal(10,2)" json:"total"`
	ExpenseDate time.Time      `gorm:"not null" json:"expense_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeSave hook to calculate total automatically
func (e *Expense) BeforeSave(tx *gorm.DB) error {
	e.Total = e.Unit * e.PerUnitCost
	return nil
}
