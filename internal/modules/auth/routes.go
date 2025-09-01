package auth

import "github.com/gofiber/fiber/v2"

func SetupAuthRoutes(app *fiber.App, handler *Handler) {
	auth := app.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)
}
