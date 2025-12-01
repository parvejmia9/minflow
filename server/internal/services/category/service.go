package category

import (
	"errors"

	"github.com/parvejmia9/minflow/server/internal/models"
	"gorm.io/gorm"
)

// Service handles category business logic
type Service struct {
	db *gorm.DB
}

// NewService creates a new category service instance
func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

// GetAll retrieves all categories for a user (default + user-specific)
func (s *Service) GetAll(userID uint) ([]models.Category, error) {
	var categories []models.Category

	// Get default categories (user_id is null) OR categories belonging to this user
	result := s.db.Where("user_id IS NULL OR user_id = ?", userID).Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

// GetByID retrieves a single category by ID
func (s *Service) GetByID(id uint) (*models.Category, error) {
	var category models.Category

	result := s.db.First(&category, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, result.Error
	}

	return &category, nil
}

// Create creates a new user-specific category
func (s *Service) Create(category *models.Category) error {
	// User-created categories are not default
	category.IsDefault = false
	result := s.db.Create(category)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// SeedDefaultCategories creates default categories if they don't exist
func (s *Service) SeedDefaultCategories() error {
	defaultCategories := []models.Category{
		{Name: "Food & Dining", IsDefault: true, UserID: nil},
		{Name: "Transportation", IsDefault: true, UserID: nil},
		{Name: "Shopping", IsDefault: true, UserID: nil},
		{Name: "Entertainment", IsDefault: true, UserID: nil},
		{Name: "Bills & Utilities", IsDefault: true, UserID: nil},
		{Name: "Healthcare", IsDefault: true, UserID: nil},
		{Name: "Education", IsDefault: true, UserID: nil},
		{Name: "Personal Care", IsDefault: true, UserID: nil},
		{Name: "Travel", IsDefault: true, UserID: nil},
		{Name: "Other", IsDefault: true, UserID: nil},
	}

	for _, cat := range defaultCategories {
		// Check if category already exists
		var existing models.Category
		result := s.db.Where("name = ? AND is_default = true", cat.Name).First(&existing)
		if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Create if doesn't exist
			if err := s.db.Create(&cat).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// Update updates an existing category
func (s *Service) Update(id uint, category *models.Category) error {
	result := s.db.Model(&models.Category{}).Where("id = ?", id).Updates(category)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}

// Delete soft deletes a category
func (s *Service) Delete(id uint) error {
	result := s.db.Delete(&models.Category{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}
