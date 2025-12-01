package services

import (
	"errors"

	"github.com/parvejmia9/minflow/server/internal/db"
	"github.com/parvejmia9/minflow/server/internal/models"
	"gorm.io/gorm"
)

// GetAllCategories retrieves all categories from the database
func GetAllCategories() ([]models.Category, error) {
	var categories []models.Category

	result := db.DB.Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

// GetCategoryByID retrieves a single category by ID
func GetCategoryByID(id uint) (*models.Category, error) {
	var category models.Category

	result := db.DB.First(&category, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, result.Error
	}

	return &category, nil
}
