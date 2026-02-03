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

	// Create map strictly with authenticated user
	fav := map[string]interface{}{
		"user_id": float64(uid.(uint)),
		"car_id":  float64(carID),
	}

	if err := h.Usecase.AddFavorite(&fav, uint(carID)); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Added to favorites"})
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
