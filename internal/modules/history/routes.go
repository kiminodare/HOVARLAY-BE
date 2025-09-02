package history

import "github.com/gofiber/fiber/v2"

const IDRoute = "/history/:id"

func SetupHistoryRoutes(router fiber.Router, handler *Handler) {
	router.Post("/history", handler.Create)
	router.Put(IDRoute, handler.Update)
	router.Get("/histories", handler.GetByUser)
	router.Get(IDRoute, handler.GetByID)
	router.Delete(IDRoute, handler.Delete)
}
