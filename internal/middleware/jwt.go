package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kiminodare/HOVARLAY-BE/internal/utils"
	"strings"
)

type JWTMiddleware struct {
	jwtUtil *utils.AESJWTUtil
}

func NewJWTMiddleware(jwtUtil *utils.AESJWTUtil) *JWTMiddleware {
	return &JWTMiddleware{jwtUtil: jwtUtil}
}

func (m *JWTMiddleware) Auth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return Error(c, "Authorization header is required", fiber.StatusUnauthorized)
	}
	token := strings.Split(authHeader, "Bearer ")
	if len(token) != 2 {
		return Error(c, "Invalid token format", fiber.StatusUnauthorized)
	}

	claims, err := m.jwtUtil.VerifyToken(token[1])
	if err != nil {
		return Error(c, "Invalid token", fiber.StatusUnauthorized)
	}

	c.Locals("user_id", claims.UserID.String())
	c.Locals("user_email", claims.Email)
	return c.Next()
}
