package http

import (
	"Backend_Go/internal/entities"
	"Backend_Go/internal/usecases/auth"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	Usecase auth.AuthUsecase
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.Usecase.RegisterUser(req.Name, req.Email, req.Phone, req.Password, req.Role); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Customer registration successful",
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	access, refresh, user, err := h.Usecase.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	response := fiber.Map{
		"access_token":  access,
		"refresh_token": refresh,
		"user_id":       user.ID,
		"role":          user.Role,
		"message":       "Login successful",
	}

	// If user is a dealer, fetch and include dealer ID
	if user.Role == "dealer" {
		var dealer entities.Dealer
		if err := h.Usecase.GetDealerByUserID(user.ID, &dealer); err == nil {
			response["dealer_id"] = dealer.ID
		}
	}

	return c.Status(200).JSON(response)
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
