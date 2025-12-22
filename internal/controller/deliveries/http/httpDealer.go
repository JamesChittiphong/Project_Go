package http

import "github.com/gofiber/fiber/v2"

type DealerHandler struct {
	Repo interface {
		Create(interface{}) error
		FindAll(interface{}) error
		FindByID(uint, interface{}) error
		Update(interface{}) error
	}
}

// POST /dealers
func (h *DealerHandler) CreateDealer(c *fiber.Ctx) error {
	var dealer map[string]interface{}
	if err := c.BodyParser(&dealer); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.Repo.Create(&dealer); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(dealer)
}

// GET /dealers
func (h *DealerHandler) GetDealers(c *fiber.Ctx) error {
	var dealers []map[string]interface{}
	if err := h.Repo.FindAll(&dealers); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(dealers)
}

// GET /dealers/:id
func (h *DealerHandler) GetDealer(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var dealer map[string]interface{}
	if err := h.Repo.FindByID(uint(id), &dealer); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "ไม่พบร้านค้า"})
	}
	return c.JSON(dealer)
}

// PUT /dealers/:id
func (h *DealerHandler) UpdateDealer(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var dealer map[string]interface{}
	if err := c.BodyParser(&dealer); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	dealer["id"] = id
	if err := h.Repo.Update(&dealer); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(dealer)
}
