package history

import "github.com/gofiber/fiber/v2"

func SetupHistoryRoutes(router fiber.Router, handler *Handler) {
	router.Post("/history", handler.Create)
	router.Put("/history/:id", handler.Update)
	router.Get("/histories", handler.GetByUser)
	router.Get("/history/:id", handler.GetByID)
	router.Delete("/history/:id", handler.Delete)
}
