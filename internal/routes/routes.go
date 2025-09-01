package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated"
	"github.com/kiminodare/HOVARLAY-BE/internal/middleware"
	"github.com/kiminodare/HOVARLAY-BE/internal/modules/auth"
	"github.com/kiminodare/HOVARLAY-BE/internal/modules/history"
	"github.com/kiminodare/HOVARLAY-BE/internal/modules/user"
	"github.com/kiminodare/HOVARLAY-BE/internal/utils"
)

func SetupRoutes(app *fiber.App, jwtMiddleware *middleware.JWTMiddleware, client *generated.Client, jwtUtil *utils.AESJWTUtil) {

	userRepository := user.NewUserRepository(client)
	userService := user.NewUserService(userRepository)

	authService := auth.NewAuthService(userService, jwtUtil)
	authHandler := auth.NewHandler(authService)

	auth.SetupAuthRoutes(app, authHandler)

	api := app.Group("/api")
	api.Use(jwtMiddleware.Auth)

	historyRepository := history.NewHistoryRepository(client)
	historyService := history.NewService(historyRepository)
	historyHandler := history.NewHandler(historyService)

	history.SetupHistoryRoutes(api, historyHandler)
}
