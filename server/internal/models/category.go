package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"size:100;not null"`
	UserID    *uint  `json:"user_id" gorm:"index"`            // null for default categories
	IsDefault bool   `json:"is_default" gorm:"default:false"` // true for system categories

	// For soft deletes and timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
