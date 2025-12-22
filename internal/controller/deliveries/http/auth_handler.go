package http

import (
	"Backend_Go/internal/usecases"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	Usecase usecases.AuthUsecase
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return h.Usecase.RegisterUser(req.Name, req.Email, req.Password, req.Role)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	access, refresh, err := h.Usecase.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	c.BodyParser(&req)
	token, err := h.Usecase.Refresh(req.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"access_token": token})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	c.BodyParser(&req)
	return h.Usecase.Logout(req.RefreshToken)
}
