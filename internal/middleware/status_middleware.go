package middleware

import (
	"Backend_Go/internal/entities"
	"Backend_Go/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

// RequireActiveDealer checks if the user is a dealer AND if the dealer status is approved
func RequireActiveDealer(dealerRepo *repositories.DealerRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. User ID from Locals
		userIDVal := c.Locals("user_id")
		if userIDVal == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: No user ID",
			})
		}
		userID, ok := userIDVal.(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid user ID",
			})
		}

		// 2. Find Dealer by UserID
		var dealer entities.Dealer
		if err := dealerRepo.FindByUserID(userID, &dealer); err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden: Dealer profile not found",
			})
		}

		// 3. Check Status
		// Fallback: check IsApproved if Status is empty (migration phase safety)
		isApproved := dealer.Status == "approved"
		if dealer.Status == "" && dealer.IsApproved {
			isApproved = true
		}

		if !isApproved {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden: Dealer account is " + dealer.Status,
			})
		}

		// Store dealer_id for convenience
		c.Locals("dealer_id", dealer.ID)

		return c.Next()
	}
}
