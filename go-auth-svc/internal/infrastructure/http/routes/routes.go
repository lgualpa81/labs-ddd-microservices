package routes

import (
	"poc-auth-svc/internal/infrastructure/http/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authHandler *handlers.AuthHandler) {
	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/validate", authHandler.ValidateToken)
}
