package http

import (
	"Backend_Go/internal/usecases"

	"github.com/gofiber/fiber/v2"
)

type LeadHandler struct {
	Usecase *usecases.LeadUsecase
}

// POST /leads
func (h *LeadHandler) CreateLead(c *fiber.Ctx) error {
	var lead map[string]interface{}
	if err := c.BodyParser(&lead); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	carID := uint(lead["car_id"].(float64))
	dealerID := uint(lead["dealer_id"].(float64))
	if err := h.Usecase.CreateLead(&lead, carID, dealerID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "ส่งข้อมูลติดต่อเรียบร้อย"})
}

// GET /dealers/:id/leads
func (h *LeadHandler) GetLeadsByDealer(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var leads []map[string]interface{}
	if err := h.Usecase.GetLeadsByDealer(uint(id), &leads); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(leads)
}
