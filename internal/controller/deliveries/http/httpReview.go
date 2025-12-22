package http

import (
	"Backend_Go/internal/usecases"

	"github.com/gofiber/fiber/v2"
)

type ReviewHandler struct {
	Usecase *usecases.ReviewUsecase
}

// POST /reviews
func (h *ReviewHandler) CreateReview(c *fiber.Ctx) error {
	var review map[string]interface{}
	if err := c.BodyParser(&review); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	dealerID := uint(review["dealer_id"].(float64))
	if err := h.Usecase.CreateReview(&review, dealerID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "รีวิวสำเร็จ"})
}

// GET /dealers/:id/reviews
func (h *ReviewHandler) GetReviewsByDealer(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var reviews []map[string]interface{}
	if err := h.Usecase.GetReviewsByDealer(uint(id), &reviews); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(reviews)
}
