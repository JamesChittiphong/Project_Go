package http

import (
	"Backend_Go/internal/entities"
	carimage "Backend_Go/internal/usecases/car_image"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CarImageHandler struct {
	Usecase *carimage.CarImageUsecase
}

// POST /cars/:id/images - AddImages creates car images (supports multiple files)
func (h *CarImageHandler) AddImages(c *fiber.Ctx) error {
	// 1️⃣ parse carID
	carIDParam := c.Params("id")
	carID, err := strconv.ParseUint(carIDParam, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid car_id",
		})
	}

	// 2️⃣ get multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid multipart form",
		})
	}

	files := form.File["images"]
	if len(files) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "at least one image file is required",
		})
	}

	// 3️⃣ สร้างโฟลเดอร์ (อิงจาก WORKDIR)
	uploadDir := filepath.Join("uploads", "cars", carIDParam)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "cannot create upload directory",
		})
	}

	var createdImages []*entities.CarImage

	// 4️⃣ loop save files
	for i, file := range files {
		// ป้องกันชื่อไฟล์ซ้ำ
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
		fullPath := filepath.Join(uploadDir, filename)

		// 🔥 สำคัญมาก
		if err := c.SaveFile(file, fullPath); err != nil {
			log.Println("SaveFile error:", err)
			continue
		}

		// 5️⃣ URL ที่ frontend ใช้
		imageURL := fmt.Sprintf(
			"/uploads/cars/%s/%s",
			carIDParam,
			filename,
		)

		image := &entities.CarImage{
			CarID:     uint(carID),
			ImageURL:  imageURL,
			SortOrder: i,
		}

		// 6️⃣ save DB
		if err := h.Usecase.CreateCarImage(image); err != nil {
			log.Println("DB save error:", err)
			continue
		}

		createdImages = append(createdImages, image)
	}

	if len(createdImages) == 0 {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to save any images",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "อัปโหลดรูปภาพสำเร็จ",
		"images":  createdImages,
		"count":   len(createdImages),
	})
}

// GET /cars/:id/images - GetImages retrieves all images for a car
func (h *CarImageHandler) GetImages(c *fiber.Ctx) error {
	carID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid car_id"})
	}

	images, err := h.Usecase.GetCarImages(uint(carID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"images": images,
	})
}

// GET /images/:id - GetImage retrieves a specific image
func (h *CarImageHandler) GetImage(c *fiber.Ctx) error {
	imageID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid image_id"})
	}

	image, err := h.Usecase.GetCarImage(uint(imageID))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "image not found"})
	}

	return c.JSON(image)
}

// PUT /images/:id - UpdateImage updates a car image
func (h *CarImageHandler) UpdateImage(c *fiber.Ctx) error {
	imageID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid image_id"})
	}

	var req struct {
		ImageURL  string `json:"image_url"`
		SortOrder int    `json:"sort_order"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	image := &entities.CarImage{
		ImageURL:  req.ImageURL,
		SortOrder: req.SortOrder,
	}
	image.ID = uint(imageID)

	if err := h.Usecase.UpdateCarImage(image); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "image updated successfully",
		"image":   image,
	})
}

// DELETE /images/:id - DeleteImage deletes a car image
func (h *CarImageHandler) DeleteImage(c *fiber.Ctx) error {
	imageID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid image_id"})
	}

	if err := h.Usecase.DeleteCarImage(uint(imageID)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "image deleted successfully",
	})
}

// DELETE /cars/:id/images - DeleteImages deletes all images for a car
func (h *CarImageHandler) DeleteImages(c *fiber.Ctx) error {
	carID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid car_id"})
	}

	if err := h.Usecase.DeleteCarImages(uint(carID)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "all images deleted successfully",
	})
}
