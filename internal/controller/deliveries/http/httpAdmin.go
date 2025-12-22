package http

import (
	"Backend_Go/internal/entities"
	"Backend_Go/internal/usecases"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	Usecase *usecases.AdminUsecase
}

// GET /admin/users
func (h *AdminHandler) GetUsers(c *fiber.Ctx) error {
	var users []map[string]interface{}
	if err := h.Usecase.GetAllUsers(&users); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

// GET /admin/dealers
func (h *AdminHandler) GetDealers(c *fiber.Ctx) error {
	var dealers []*entities.Dealer
	if err := h.Usecase.GetAllDealers(&dealers); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(dealers)
}

// GET /admin/reports
func (h *AdminHandler) GetReports(c *fiber.Ctx) error {
	var reports []map[string]interface{}
	if err := h.Usecase.GetAllReports(&reports); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(reports)
}

// POST /admin/dealers/:id/approve
func (h *AdminHandler) ApproveDealer(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var req struct {
		Approve bool `json:"approve"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.Usecase.SetDealerApproval(uint(id), req.Approve); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "ดำเนินการเรียบร้อย"})
}

// POST /admin/cars/:id/hide
func (h *AdminHandler) HideCar(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var req struct {
		Hide bool `json:"hide"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.Usecase.SetCarHidden(uint(id), req.Hide); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "ดำเนินการเรียบร้อย"})
}

// POST /admin/cars/:id/flag
func (h *AdminHandler) FlagCar(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.Usecase.FlagCar(uint(id), req.Reason); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "ดำเนินการเรียบร้อย"})
}

// DELETE /admin/cars/:id
func (h *AdminHandler) DeleteCar(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := h.Usecase.DeleteCar(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "ลบรถเรียบร้อยโดยแอดมิน"})
}

func interfaceSlice[T any](in []T) []interface{} {
	out := make([]interface{}, len(in))
	for i := range in {
		out[i] = in[i]
	}
	return out
}
