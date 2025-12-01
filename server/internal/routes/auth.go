package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/parvejmia9/minflow/server/internal/handlers"
)

func SetupAuthRoutes(router fiber.Router, authHandler *handlers.AuthHandler) {
	// POST /auth/signup - Register new user
	router.Post("/auth/signup", authHandler.Signup)

	// POST /auth/login - Login user
	router.Post("/auth/login", authHandler.Login)
}
