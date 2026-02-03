package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
)

// RequireRole returns a Fiber middleware that validates JWT and checks the role claim.
func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(401).JSON(fiber.Map{"error": "session expired"})
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(401).JSON(fiber.Map{"error": "session expired"})
		}
		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "session expired"})
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "session expired"})
		}
		r, _ := claims["role"].(string)
		idf, _ := claims["user_id"].(float64)

		if r != role {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}

		// store in locals for handlers
		c.Locals("user_id", uint(idf))
		c.Locals("role", r)
		return c.Next()
	}
}

// RequireAuth validates JWT and sets user info, without checking role.
func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(401).JSON(fiber.Map{"error": "session expired"})
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(401).JSON(fiber.Map{"error": "session expired"})
		}
		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "session expired"})
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "session expired"})
		}
		idf, _ := claims["user_id"].(float64)
		r, _ := claims["role"].(string)
		c.Locals("user_id", uint(idf))
		c.Locals("role", r)
		return c.Next()
	}
}
