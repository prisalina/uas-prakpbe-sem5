package middleware

import (
	"strings"
	"uas-pbe-praksem5/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Missing Authorization header"})
		}
		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid Authorization header"})
		}
		tokenStr := parts[1]
		claims, err := utils.ValidateAccessToken(tokenStr)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token", "detail": err.Error()})
		}
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role_name", claims.RoleName)
		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		r := c.Locals("role_name")
		if r == nil || r.(string) != "Admin" {
			return c.Status(403).JSON(fiber.Map{"error": "Only Admin allowed"})
		}
		return c.Next()
	}
}

func RoleIs(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		r := c.Locals("role_name")
		if r == nil || r.(string) != role {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
		return c.Next()
	}
}
