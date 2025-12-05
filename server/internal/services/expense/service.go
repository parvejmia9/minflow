package expense

import (
	"errors"
	"time"

	"github.com/parvejmia9/minflow/server/internal/models"
	"gorm.io/gorm"
)

// Service handles expense business logic
type Service struct {
	db *gorm.DB
}

// NewService creates a new expense service instance
func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

// CreateExpenseInput represents the input for creating an expense
type CreateExpenseInput struct {
	Name        string    `json:"name" validate:"required"`
	CategoryID  uint      `json:"category_id" validate:"required"`
	Unit        float64   `json:"unit" validate:"required,gt=0"`
	PerUnitCost float64   `json:"per_unit_cost" validate:"required,gt=0"`
	ExpenseDate time.Time `json:"expense_date"`
}

// AnalyticsQuery represents the query parameters for analytics
type AnalyticsQuery struct {
	StartDate time.Time
	EndDate   time.Time
	UserID    uint
}

// AnalyticsResult represents the analytics data
type AnalyticsResult struct {
	TotalExpenses     float64           `json:"total_expenses"`
	ExpenseCount      int64             `json:"expense_count"`
	ByCategory        []CategoryExpense `json:"by_category"`
	DailyExpenses     []DailyExpense    `json:"daily_expenses"`
	AverageDailySpend float64           `json:"average_daily_spend"`
	DateRange         DateRange         `json:"date_range"`
}

type CategoryExpense struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Total        float64 `json:"total"`
	Count        int64   `json:"count"`
}

type DailyExpense struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}

type DateRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Create creates a new expense
func (s *Service) Create(userID uint, input CreateExpenseInput) (*models.Expense, error) {
	// Verify category exists
	var category models.Category
	if err := s.db.First(&category, input.CategoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	// Set expense date to now if not provided
	if input.ExpenseDate.IsZero() {
		input.ExpenseDate = time.Now()
	}

	expense := &models.Expense{
		Name:        input.Name,
		CategoryID:  input.CategoryID,
		UserID:      userID,
		Unit:        input.Unit,
		PerUnitCost: input.PerUnitCost,
		ExpenseDate: input.ExpenseDate,
	}

	// Total is calculated automatically in BeforeSave hook
	if err := s.db.Create(expense).Error; err != nil {
		return nil, err
	}

	// Load category relationship
	s.db.Preload("Category").First(expense, expense.ID)

	return expense, nil
}

// GetByUser retrieves all expenses for a user
func (s *Service) GetByUser(userID uint, limit, offset int) ([]models.Expense, int64, error) {
	var expenses []models.Expense
	var total int64

	// Count total
	s.db.Model(&models.Expense{}).Where("user_id = ?", userID).Count(&total)

	// Get paginated results
	err := s.db.
		Preload("Category").
		Where("user_id = ?", userID).
		Order("expense_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&expenses).Error

	if err != nil {
		return nil, 0, err
	}

	return expenses, total, nil
}

// GetByID retrieves a single expense by ID
func (s *Service) GetByID(id, userID uint) (*models.Expense, error) {
	var expense models.Expense

	err := s.db.
		Preload("Category").
		Where("id = ? AND user_id = ?", id, userID).
		First(&expense).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("expense not found")
		}
		return nil, err
	}

	return &expense, nil
}

// GetDateRange gets the first and last expense dates for a user
func (s *Service) GetDateRange(userID uint) (*DateRange, error) {
	var firstExpense, lastExpense models.Expense

	// Get first expense
	err := s.db.Where("user_id = ?", userID).
		Order("expense_date ASC").
		First(&firstExpense).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no expenses found")
		}
		return nil, err
	}

	// Get last expense (or use today)
	err = s.db.Where("user_id = ?", userID).
		Order("expense_date DESC").
		First(&lastExpense).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	endDate := lastExpense.ExpenseDate
	today := time.Now()
	if endDate.After(today) {
		endDate = today
	}

	return &DateRange{
		Start: firstExpense.ExpenseDate,
		End:   endDate,
	}, nil
}

// GetAnalytics generates analytics for expenses within a date range
func (s *Service) GetAnalytics(query AnalyticsQuery) (*AnalyticsResult, error) {
	result := &AnalyticsResult{
		DateRange: DateRange{
			Start: query.StartDate,
			End:   query.EndDate,
		},
	}

	// Get total expenses and count
	var totalSum struct {
		Total float64
		Count int64
	}

	err := s.db.Model(&models.Expense{}).
		Select("COALESCE(SUM(total), 0) as total, COUNT(*) as count").
		Where("user_id = ? AND expense_date BETWEEN ? AND ?", query.UserID, query.StartDate, query.EndDate).
		Scan(&totalSum).Error

	if err != nil {
		return nil, err
	}

	result.TotalExpenses = totalSum.Total
	result.ExpenseCount = totalSum.Count

	// Get expenses by category
	err = s.db.Model(&models.Expense{}).
		Select("categories.id as category_id, categories.name as category_name, COALESCE(SUM(expenses.total), 0) as total, COUNT(expenses.id) as count").
		Joins("LEFT JOIN categories ON categories.id = expenses.category_id").
		Where("expenses.user_id = ? AND expenses.expense_date BETWEEN ? AND ?", query.UserID, query.StartDate, query.EndDate).
		Group("categories.id, categories.name").
		Order("total DESC").
		Scan(&result.ByCategory).Error

	if err != nil {
		return nil, err
	}

	// Get daily expenses
	err = s.db.Model(&models.Expense{}).
		Select("DATE(expense_date) as date, COALESCE(SUM(total), 0) as total").
		Where("user_id = ? AND expense_date BETWEEN ? AND ?", query.UserID, query.StartDate, query.EndDate).
		Group("DATE(expense_date)").
		Order("date ASC").
		Scan(&result.DailyExpenses).Error

	if err != nil {
		return nil, err
	}

	// Calculate average daily spend
	days := query.EndDate.Sub(query.StartDate).Hours() / 24
	if days > 0 && result.TotalExpenses > 0 {
		result.AverageDailySpend = result.TotalExpenses / days
	}

	return result, nil
}

// Delete soft deletes an expense
func (s *Service) Delete(id, userID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Expense{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("expense not found")
	}
	return nil
}
