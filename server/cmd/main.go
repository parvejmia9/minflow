package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/parvejmia9/minflow/server/internal/db"
	"github.com/parvejmia9/minflow/server/internal/handlers"
	"github.com/parvejmia9/minflow/server/internal/models"
	"github.com/parvejmia9/minflow/server/internal/routes"
	"github.com/parvejmia9/minflow/server/internal/services/auth"
	"github.com/parvejmia9/minflow/server/internal/services/category"
	"github.com/parvejmia9/minflow/server/internal/services/expense"
	"github.com/parvejmia9/minflow/server/internal/services/user"
)

func main() {
	// Load .env file
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: .env file not found, using environment variables or defaults")
	}

	// Connect to database
	db.ConnectDB()

	// Auto migrate database models
	err := db.DB.AutoMigrate(&models.User{}, &models.Category{}, &models.Expense{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed default categories
	tempCategoryService := category.NewService(db.DB)
	if err := tempCategoryService.SeedDefaultCategories(); err != nil {
		log.Println("Warning: Failed to seed default categories:", err)
	} else {
		log.Println("Default categories seeded successfully")
	}

	// Get JWT secret from environment or use default (change in production!)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production!")
	}

	// Initialize services with dependency injection
	authService := auth.NewService(db.DB, jwtSecret)
	categoryService := category.NewService(db.DB)
	expenseService := expense.NewService(db.DB)
	userService := user.NewService(db.DB)

	// Initialize handlers with service dependencies
	authHandler := handlers.NewAuthHandler(authService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	expenseHandler := handlers.NewExpenseHandler(expenseService)
	userHandler := handlers.NewUserHandler(userService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Setup routes with handler dependencies
	routes.SetupRoutes(app, authService, authHandler, categoryHandler, expenseHandler, userHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	log.Printf("Server starting on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
