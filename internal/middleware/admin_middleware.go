package middleware

import (
	"Backend_Go/internal/usecases/admin"

	"github.com/gofiber/fiber/v2"
)

// AdminOnly returns a middleware that checks if the authenticated user is an admin
// It relies on "user_id" being present in locals (set by JWTAuth middleware)
func AdminOnly(usecase *admin.AdminUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Get user_id from Locals (set by JWT middleware)
		userID := c.Locals("user_id")
		if userID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: No user ID found in session",
			})
		}

		uid, ok := userID.(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid user ID format",
			})
		}

		// 2. Check logic via Usecase
		isAdmin, err := usecase.IsAdmin(uid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error checking permissions",
			})
		}

		if !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden: Admin access required",
			})
		}

		// 3. User is admin, proceed
		return c.Next()
	}
}
