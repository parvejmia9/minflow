package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/parvejmia9/minflow/server/internal/handlers"
	"github.com/parvejmia9/minflow/server/internal/middleware"
	"github.com/parvejmia9/minflow/server/internal/services/auth"
)

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

func SetupRoutes(
	app *fiber.App,
	authService *auth.Service,
	authHandler *handlers.AuthHandler,
	categoryHandler *handlers.CategoryHandler,
	expenseHandler *handlers.ExpenseHandler,
	userHandler *handlers.UserHandler,
) {
	api := app.Group("/api")

	// Health check (public)
	api.Get("/health", healthCheck)

	// Auth routes (public)
	SetupAuthRoutes(api, authHandler)

	// Protected routes (require authentication)
	protected := api.Group("", middleware.AuthMiddleware(authService))

	// Category routes
	SetupCategoryRoutes(protected, categoryHandler)

	// Expense routes
	SetupExpenseRoutes(protected, expenseHandler)

	// User routes (includes both user and admin routes)
	SetupUserRoutes(protected, userHandler)
}
