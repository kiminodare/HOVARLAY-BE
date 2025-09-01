package middleware

import "github.com/gofiber/fiber/v2"

type ApiResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Error      interface{} `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

func Success(c *fiber.Ctx, data interface{}, message string, pagination *Pagination) error {
	return c.JSON(ApiResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

func Error(c *fiber.Ctx, message string, statusCode int) error {
	c.Status(statusCode)
	return c.JSON(ApiResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

func ValidationError(c *fiber.Ctx, errors []string) error {
	c.Status(400)
	return c.JSON(ApiResponse{
		Success: false,
		Message: "Validation failed",
		Error:   errors,
		Data:    nil,
	})
}
