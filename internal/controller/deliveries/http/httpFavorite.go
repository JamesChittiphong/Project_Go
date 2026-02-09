package http

import (
	"Backend_Go/internal/entities"
	"Backend_Go/internal/usecases/favorite"

	"github.com/gofiber/fiber/v2"
)

type FavoriteHandler struct {
	Usecase *favorite.FavoriteUsecase
}

// POST /favorites
// POST /favorites
// POST /favorites
func (h *FavoriteHandler) AddFavorite(c *fiber.Ctx) error {
	var body struct {
		UserID uint `json:"user_id"`
		CarID  uint `json:"car_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	status, err := h.Usecase.ToggleFavorite(body.UserID, body.CarID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": status, "message": "Favorite toggled successfully"})
}

// GET /users/:id/favorites
func (h *FavoriteHandler) GetFavoritesByUser(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var favs []entities.Favorite
	if err := h.Usecase.GetFavoritesByUser(uint(id), &favs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// Transform to return just cars? Or return Favorite objects?
	// Returning Favorite objects is fine, frontend can extract car.
	return c.JSON(favs)
}

// DELETE /users/:id/favorites/:car_id
func (h *FavoriteHandler) RemoveFavorite(c *fiber.Ctx) error {
	userID, _ := c.ParamsInt("id")
	carID, _ := c.ParamsInt("car_id")

	if err := h.Usecase.RemoveFavorite(uint(userID), uint(carID)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "ลบรายการโปรดเรียบร้อย"})
}

// Secure Favorites
// GET /favorites/me
func (h *FavoriteHandler) GetMyFavorites(c *fiber.Ctx) error {
	uid := c.Locals("user_id")
	if uid == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	var favs []entities.Favorite
	if err := h.Usecase.GetFavoritesByUser(uid.(uint), &favs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(favs)
}

// POST /favorites/:car_id
func (h *FavoriteHandler) AddFavoriteMe(c *fiber.Ctx) error {
	uid := c.Locals("user_id")
	if uid == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	carID, err := c.ParamsInt("car_id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid car id"})
	}

	userID := uid.(uint)

	status, err := h.Usecase.ToggleFavorite(userID, uint(carID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": status, "message": "Favorite toggled successfully"})
}

// DELETE /favorites/:car_id
func (h *FavoriteHandler) RemoveFavoriteMe(c *fiber.Ctx) error {
	uid := c.Locals("user_id")
	if uid == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	carID, err := c.ParamsInt("car_id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid car id"})
	}

	if err := h.Usecase.RemoveFavorite(uid.(uint), uint(carID)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Removed from favorites"})
}
