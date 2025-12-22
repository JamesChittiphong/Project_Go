package http

import (
	"Backend_Go/internal/entities"
	"Backend_Go/internal/usecases"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	Usecase *usecases.UserUsecase
}

// POST /users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user entities.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.Usecase.CreateUser(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

// GET /users/:id
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	user, err := h.Usecase.GetUserByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "ไม่พบผู้ใช้"})
	}
	return c.JSON(user)
}

// PUT /users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var user entities.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	user.ID = uint(id)
	if err := h.Usecase.UpdateUser(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

// DELETE /users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := h.Usecase.DeleteUser(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "ลบผู้ใช้แล้ว"})
}
