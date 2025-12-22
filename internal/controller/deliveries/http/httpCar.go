package http

import (
	"Backend_Go/internal/entities"
	"Backend_Go/internal/usecases"
	"Backend_Go/utils"

	"github.com/gofiber/fiber/v2"
)

type CarHandler struct {
	Usecase *usecases.CarUsecase
}

// POST /cars
func (h *CarHandler) CreateCar(c *fiber.Ctx) error {
	var req struct {
		Car    entities.Car        `json:"car"`
		Images []entities.CarImage `json:"images"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	images := make([]interface{}, len(req.Images))
	for i := range req.Images {
		images[i] = &req.Images[i]
	}
	if err := h.Usecase.CreateCar(&req.Car, images); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "สร้างรถสำเร็จ"})
}

// GET /cars
func (h *CarHandler) GetCars(c *fiber.Ctx) error {
	var cars []*entities.Car
	if err := h.Usecase.GetAllCars(&cars); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cars)
}

// GET /cars/:id
func (h *CarHandler) GetCarDetail(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var car entities.Car
	var images []entities.CarImage
	if err := h.Usecase.GetCarDetail(uint(id), &car, &images); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"car": car, "images": images})
}

// DELETE /cars/:id
func (h *CarHandler) DeleteCar(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	// owner check
	uid := c.Locals("user_id")
	if uid == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	userID := uid.(uint)
	var dealer entities.Dealer
	if err := h.Usecase.DealerRepo.FindByUserID(userID, &dealer); err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}
	var car entities.Car
	if err := h.Usecase.CarRepo.FindByID(uint(id), &car); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	if car.DealerID != dealer.ID {
		return c.Status(403).JSON(fiber.Map{"error": "not owner of car"})
	}
	if err := h.Usecase.DeleteCar(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "ลบรถเรียบร้อย"})
}

// PUT /cars/:id
func (h *CarHandler) UpdateCar(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var car entities.Car
	if err := c.BodyParser(&car); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	car.ID = uint(id)
	// owner check
	uid := c.Locals("user_id")
	if uid == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	userID := uid.(uint)
	var dealer entities.Dealer
	if err := h.Usecase.DealerRepo.FindByUserID(userID, &dealer); err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}
	var existing entities.Car
	if err := h.Usecase.CarRepo.FindByID(uint(id), &existing); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	if existing.DealerID != dealer.ID {
		return c.Status(403).JSON(fiber.Map{"error": "not owner of car"})
	}
	if err := h.Usecase.UpdateCar(&car); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "แก้ไขรถเรียบร้อย"})
}

// PATCH /cars/:id/status
func (h *CarHandler) SetStatus(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// owner check
	uid := c.Locals("user_id")
	if uid == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	userID := uid.(uint)
	var dealer entities.Dealer
	if err := h.Usecase.DealerRepo.FindByUserID(userID, &dealer); err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}
	var car entities.Car
	if err := h.Usecase.CarRepo.FindByID(uint(id), &car); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	if car.DealerID != dealer.ID {
		return c.Status(403).JSON(fiber.Map{"error": "not owner"})
	}
	if err := h.Usecase.SetStatus(uint(id), req.Status); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "อัปเดตสถานะเรียบร้อย"})
}

// POST /cars/:id/contact
func (h *CarHandler) RecordContact(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var req struct {
		Via      string `json:"via"`
		DealerID uint   `json:"dealer_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.Usecase.RecordContact(uint(id), req.DealerID, req.Via); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// async notification webhook
	go func() {
		webhook := utils.GetEnv("NOTIFY_WEBHOOK_URL", "")
		if webhook == "" {
			return
		}
		payload := map[string]interface{}{
			"car_id":    id,
			"dealer_id": req.DealerID,
			"via":       req.Via,
		}
		_ = utils.SendWebhookNotification(webhook, payload)
	}()

	return c.JSON(fiber.Map{"message": "บันทึกการติดต่อเรียบร้อย"})
}

// POST /cars/:id/promote
func (h *CarHandler) PromoteCar(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var req struct {
		Days int `json:"days"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if req.Days <= 0 {
		req.Days = 7
	}
	// owner check
	uid := c.Locals("user_id")
	if uid == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	userID := uid.(uint)
	var dealer entities.Dealer
	if err := h.Usecase.DealerRepo.FindByUserID(userID, &dealer); err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}
	var car entities.Car
	if err := h.Usecase.CarRepo.FindByID(uint(id), &car); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	if car.DealerID != dealer.ID {
		return c.Status(403).JSON(fiber.Map{"error": "not owner"})
	}
	if err := h.Usecase.PromoteCar(uint(id), req.Days); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "โปรโมทรถเรียบร้อย"})
}

// GET /cars/:id/stats
func (h *CarHandler) GetStats(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	car, err := h.Usecase.GetStats(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"call_count":     car.CallCount,
		"line_count":     car.LineCount,
		"lead_count":     car.LeadCount,
		"is_promoted":    car.IsPromoted,
		"promoted_until": car.PromotedUntil,
	})
}
