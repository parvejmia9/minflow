package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/parvejmia9/minflow/server/internal/services/expense"
)

// ExpenseHandler handles HTTP requests for expenses
type ExpenseHandler struct {
	expenseService *expense.Service
}

// NewExpenseHandler creates a new expense handler
func NewExpenseHandler(expenseService *expense.Service) *ExpenseHandler {
	return &ExpenseHandler{
		expenseService: expenseService,
	}
}

// Create handles POST /expenses
func (h *ExpenseHandler) Create(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(uint)

	var input expense.CreateExpenseInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Validate required fields
	if input.Name == "" || input.CategoryID == 0 || input.Unit <= 0 || input.PerUnitCost <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Name, category_id, unit, and per_unit_cost are required and must be positive",
		})
	}

	expense, err := h.expenseService.Create(userID, input)
	if err != nil {
		if err.Error() == "category not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to create expense",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    expense,
	})
}

// GetAll handles GET /expenses
func (h *ExpenseHandler) GetAll(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	// Parse pagination params
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	expenses, total, err := h.expenseService.GetByUser(userID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch expenses",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    expenses,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

// GetByID handles GET /expenses/:id
func (h *ExpenseHandler) GetByID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid expense ID",
		})
	}

	expense, err := h.expenseService.GetByID(uint(id), userID)
	if err != nil {
		if err.Error() == "expense not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch expense",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    expense,
	})
}

// GetDateRange handles GET /expenses/date-range
func (h *ExpenseHandler) GetDateRange(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	dateRange, err := h.expenseService.GetDateRange(userID)
	if err != nil {
		if err.Error() == "no expenses found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch date range",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    dateRange,
	})
}

// GetAnalytics handles GET /expenses/analytics
func (h *ExpenseHandler) GetAnalytics(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	// Parse date parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "start_date and end_date are required (format: YYYY-MM-DD)",
		})
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid start_date format (use YYYY-MM-DD)",
		})
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid end_date format (use YYYY-MM-DD)",
		})
	}

	// Set to end of day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// Validate dates
	if endDate.Before(startDate) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "end_date must be after start_date",
		})
	}

	// Don't allow future dates
	if endDate.After(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "end_date cannot be in the future",
		})
	}

	query := expense.AnalyticsQuery{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	analytics, err := h.expenseService.GetAnalytics(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to generate analytics",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    analytics,
	})
}

// Delete handles DELETE /expenses/:id
func (h *ExpenseHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid expense ID",
		})
	}

	err = h.expenseService.Delete(uint(id), userID)
	if err != nil {
		if err.Error() == "expense not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to delete expense",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Expense deleted successfully",
	})
}
