package middleware

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func CheckAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Authorization header is missing",
		})
	}
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid authorization header format",
		})
	}
	token := splitToken[1]
	c.Locals("token", token)
	return c.Next()
}
