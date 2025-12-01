package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/parvejmia9/minflow/server/internal/handlers"
	"github.com/parvejmia9/minflow/server/internal/middleware"
)

func SetupUserRoutes(router fiber.Router, userHandler *handlers.UserHandler) {
	// GET /users/me - Get current user (requires auth)
	router.Get("/users/me", userHandler.GetMe)

	// Admin routes (require admin access)
	admin := router.Group("", middleware.AdminMiddleware())

	// GET /users - Get all users
	admin.Get("/users", userHandler.GetAll)

	// GET /users/:id - Get user by ID
	admin.Get("/users/:id", userHandler.GetByID)

	// GET /users/:id/stats - Get user stats
	admin.Get("/users/:id/stats", userHandler.GetStats)

	// DELETE /users/:id - Delete user
	admin.Delete("/users/:id", userHandler.Delete)
}
