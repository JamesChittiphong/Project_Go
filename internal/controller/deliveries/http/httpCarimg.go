package http

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CarImageHandler struct {
	Repo interface {
		Create(interface{}) error
		FindByCarID(uint, interface{}) error
		Delete(uint) error
	}
}

// POST /cars/:id/images  (multipart/form-data with field name `images`)
func (h *CarImageHandler) AddImage(c *fiber.Ctx) error {
	carID, _ := c.ParamsInt("id")

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid multipart form"})
	}
	files := form.File["images"]
	if len(files) == 0 {
		// try single file field
		f, ferr := c.FormFile("image")
		if ferr != nil {
			return c.Status(400).JSON(fiber.Map{"error": "no images uploaded"})
		}
		files = append(files, f)
	}

	saved := make([]map[string]interface{}, 0, len(files))
	destDir := filepath.Join("uploads", "cars", fmt.Sprintf("%d", carID))
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	for _, fh := range files {
		// create unique filename
		name := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fh.Filename)
		dest := filepath.Join(destDir, name)
		if err := c.SaveFile(fh, dest); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		rec := map[string]interface{}{
			"car_id":    carID,
			"image_url": dest,
		}
		if err := h.Repo.Create(&rec); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		saved = append(saved, rec)
	}
	return c.JSON(saved)
}

// GET /cars/:id/images
func (h *CarImageHandler) GetImages(c *fiber.Ctx) error {
	carID, _ := c.ParamsInt("id")
	var imgs []map[string]interface{}
	if err := h.Repo.FindByCarID(uint(carID), &imgs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(imgs)
}

// DELETE /images/:id
func (h *CarImageHandler) DeleteImage(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := h.Repo.Delete(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "ลบรูปเรียบร้อย"})
}
