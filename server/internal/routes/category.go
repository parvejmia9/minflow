package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/parvejmia9/minflow/server/internal/handlers"
)

func SetupCategoryRoutes(router fiber.Router, categoryHandler *handlers.CategoryHandler) {
	// GET /categories - Get all categories
	router.Get("/categories", categoryHandler.GetAll)

	// GET /categories/:id - Get single category by ID
	router.Get("/categories/:id", categoryHandler.GetByID)

	// POST /categories - Create new category
	router.Post("/categories", categoryHandler.Create)
}
