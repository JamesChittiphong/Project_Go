package http

import "github.com/gofiber/fiber/v2"

// สมัครร้านค้า (Dealer Register)
// ใช้ AuthHandler ตัวเดิม แต่บังคับ role = "dealer"
func (h *AuthHandler) RegisterDealer(c *fiber.Ctx) error {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`

		ShopName string `json:"shop_name"`
		Phone    string `json:"phone"`
		LineID   string `json:"line_id"`
		Address  string `json:"address"`
		Province string `json:"province"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.Usecase.RegisterDealer(
		req.Name,
		req.Email,
		req.Password,
		req.ShopName,
		req.Phone,
		req.LineID,
	)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "สมัครร้านค้าสำเร็จ"})
}
