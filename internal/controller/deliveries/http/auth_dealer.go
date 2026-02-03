package http

import "github.com/gofiber/fiber/v2"

// สมัครร้านค้า (Dealer Register)
func (h *AuthHandler) RegisterDealer(c *fiber.Ctx) error {
	var req struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Password  string `json:"password"`
		ShopName  string `json:"shop_name"`
		LineID    string `json:"line_id"`
		Address   string `json:"address"`
		Province  string `json:"province"`
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// validate required fields
	if req.Name == "" ||
		req.Email == "" ||
		req.Phone == "" ||
		req.Password == "" ||
		req.ShopName == "" ||
		req.Address == "" ||
		req.Province == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	err := h.Usecase.RegisterDealer(
		req.Name,
		req.Email,
		req.Phone,
		req.Password,
		req.ShopName,
		req.LineID,
		req.Address,
		req.Province,
		req.Latitude,
		req.Longitude,
	)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Dealer registration successful, waiting for approval",
	})
}
