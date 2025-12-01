package user

import (
	"errors"

	"github.com/parvejmia9/minflow/server/internal/models"
	"gorm.io/gorm"
)

// Service handles user business logic
type Service struct {
	db *gorm.DB
}

// NewService creates a new user service instance
func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

// GetAll retrieves all users (admin only)
func (s *Service) GetAll(limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Count total
	s.db.Model(&models.User{}).Count(&total)

	// Get paginated results
	err := s.db.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetByID retrieves a single user by ID
func (s *Service) GetByID(id uint) (*models.User, error) {
	var user models.User

	err := s.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (s *Service) GetByEmail(email string) (*models.User, error) {
	var user models.User

	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// Delete soft deletes a user (admin only)
func (s *Service) Delete(id uint) error {
	// Don't allow deleting admin users
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if user.IsAdmin {
		return errors.New("cannot delete admin user")
	}

	result := s.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// GetUserStats returns statistics for a user
func (s *Service) GetUserStats(userID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count total expenses
	var expenseCount int64
	s.db.Model(&models.Expense{}).Where("user_id = ?", userID).Count(&expenseCount)
	stats["total_expenses"] = expenseCount

	// Sum total spending
	var totalSpending float64
	s.db.Model(&models.Expense{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(total), 0)").
		Scan(&totalSpending)
	stats["total_spending"] = totalSpending

	// Count categories used
	var categoryCount int64
	s.db.Model(&models.Expense{}).
		Where("user_id = ?", userID).
		Distinct("category_id").
		Count(&categoryCount)
	stats["categories_used"] = categoryCount

	return stats, nil
}
