package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/parvejmia9/minflow/server/internal/handlers"
)

func SetupExpenseRoutes(router fiber.Router, expenseHandler *handlers.ExpenseHandler) {
	// All expense routes require authentication
	// POST /expenses - Create new expense
	router.Post("/expenses", expenseHandler.Create)

	// GET /expenses - Get all expenses for user (paginated)
	router.Get("/expenses", expenseHandler.GetAll)

	// GET /expenses/date-range - Get date range of expenses
	router.Get("/expenses/date-range", expenseHandler.GetDateRange)

	// GET /expenses/analytics - Get analytics
	router.Get("/expenses/analytics", expenseHandler.GetAnalytics)

	// GET /expenses/:id - Get single expense
	router.Get("/expenses/:id", expenseHandler.GetByID)

	// DELETE /expenses/:id - Delete expense
	router.Delete("/expenses/:id", expenseHandler.Delete)
}
