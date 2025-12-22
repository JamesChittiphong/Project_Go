package http

import (
	"Backend_Go/internal/usecases"

	"github.com/gofiber/fiber/v2"
)

type FavoriteHandler struct {
	Usecase *usecases.FavoriteUsecase
}

// POST /favorites
func (h *FavoriteHandler) AddFavorite(c *fiber.Ctx) error {
	var fav map[string]interface{}
	if err := c.BodyParser(&fav); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	carID := uint(fav["car_id"].(float64))
	if err := h.Usecase.AddFavorite(&fav, carID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "เพิ่มรายการโปรดแล้ว"})
}

// GET /users/:id/favorites
func (h *FavoriteHandler) GetFavoritesByUser(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var favs []map[string]interface{}
	if err := h.Usecase.GetFavoritesByUser(uint(id), &favs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(favs)
}
