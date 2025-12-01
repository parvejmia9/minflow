package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/parvejmia9/minflow/server/internal/models"
	"github.com/parvejmia9/minflow/server/internal/services/category"
)

// CategoryHandler handles HTTP requests for categories
type CategoryHandler struct {
	categoryService *category.Service
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryService *category.Service) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// GetAll handles GET /categories
func (h *CategoryHandler) GetAll(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(uint)

	categories, err := h.categoryService.GetAll(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch categories",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    categories,
		"count":   len(categories),
	})
}

// GetByID handles GET /categories/:id
func (h *CategoryHandler) GetByID(c *fiber.Ctx) error {
	// Get category ID from URL params
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid category ID",
		})
	}

	category, err := h.categoryService.GetByID(uint(id))
	if err != nil {
		if err.Error() == "category not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Category not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch category",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    category,
	})
}

// Create handles POST /categories
func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(uint)

	var category models.Category

	// Parse request body
	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Set user ID for user-specific category
	category.UserID = &userID

	// Create category
	if err := h.categoryService.Create(&category); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to create category",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    category,
	})
}
