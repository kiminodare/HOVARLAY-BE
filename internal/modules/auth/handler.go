package auth

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kiminodare/HOVARLAY-BE/internal/middleware"
	dtoAuth "github.com/kiminodare/HOVARLAY-BE/internal/modules/auth/dto"
	dtoUser "github.com/kiminodare/HOVARLAY-BE/internal/modules/user/dto"
	"github.com/kiminodare/HOVARLAY-BE/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req dtoAuth.Request
	if err := c.BodyParser(&req); err != nil {
		fmt.Printf("Body parse error: %v\n", err)
		return middleware.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	res, err := h.service.Login(c.Context(), &req)
	if err != nil {
		if errors.Is(err, utils.ErrInvalidCredentials) {
			return middleware.Error(c, "Invalid email or password", fiber.StatusUnauthorized)
		}
		return middleware.Error(c, "Failed to login", fiber.StatusInternalServerError)
	}

	c.ClearCookie("token")

	c.Cookie(&fiber.Cookie{
		Name:     "login_test",
		Value:    "simple_value",
		Path:     "/",
		MaxAge:   86400,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	return middleware.Success(c, res, "Login successful", nil)
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var req dtoUser.Request
	if err := c.BodyParser(&req); err != nil {
		return middleware.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	if err := req.Validate(); err != nil {
		return middleware.ValidationError(c, utils.FormatValidationErrors(err))
	}

	_, err := h.service.Register(c.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrEmailAlreadyExists):
			return middleware.Error(c, "Email already registered", fiber.StatusConflict)
		case errors.Is(err, utils.ErrInvalidData):
			return middleware.Error(c, "Invalid data", fiber.StatusBadRequest)
		default:
			return middleware.Error(c, "Failed to register", fiber.StatusInternalServerError)
		}
	}

	return middleware.Success(c, nil, "Registration successful", nil)
}
