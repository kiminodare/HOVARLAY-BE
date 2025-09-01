package main

import (
	"context"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated"
	"github.com/kiminodare/HOVARLAY-BE/internal/middleware"
	"github.com/kiminodare/HOVARLAY-BE/internal/routes"
	"github.com/kiminodare/HOVARLAY-BE/internal/utils"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/kiminodare/HOVARLAY-BE/internal/db"
)

func main() {
	// load env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ no .env file found, using system env")
	}

	// init DB
	client := db.NewClient()
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("❌ failed closing connection: %v", err)
		}
	}(client)

	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	// init Fiber
	app := fiber.New()
	app.Use(cors.New(
		cors.Config{
			AllowOriginsFunc: func(origin string) bool {
				for _, o := range allowedOrigins {
					if strings.TrimSpace(o) == origin {
						return true
					}
				}
				return false
			},
			AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
			AllowHeaders:     "Origin,Content-Type,Authorization,Accept",
			AllowCredentials: true,
		},
	))

	jwtUtils := utils.NewAESJWTUtil(os.Getenv("JWT_SECRET"), os.Getenv("AES_KEY"))
	jwtMiddleware := middleware.NewJWTMiddleware(jwtUtils)
	routes.SetupRoutes(app, jwtMiddleware, client, jwtUtils)

	// health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return middleware.Success(c, nil, "OK", nil)
	})

	app.Get("/test-cookie", func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:     "test",
			Value:    "testvalue",
			HTTPOnly: false,
			Secure:   false,
			SameSite: "Lax",
			Path:     "/",
		})
		return c.SendString("Cookie set")
	})

	app.Get("*", func(c *fiber.Ctx) error {
		return middleware.Error(c, "What you're looking for is not here", fiber.StatusNotFound)
	})

	// graceful shutdown
	go func() {
		if err := app.Listen(":" + os.Getenv("SERVER_PORT")); err != nil {
			log.Fatalf("❌ server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Shutdown(); err != nil {
		log.Fatalf("❌ failed to shutdown: %v", err)
	}
}
