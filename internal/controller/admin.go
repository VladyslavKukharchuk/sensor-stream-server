//nolint:wrapcheck
package controller

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"sensor-stream-server/internal/model"
)

type MeasurementService interface {
	List(ctx context.Context) ([]*model.Measurement, error)
}

type AdminController struct {
	ms MeasurementService // тільки для сторінок з вимірами
}

func NewAdminController(ms MeasurementService) *AdminController {
	return &AdminController{
		ms: ms,
	}
}

func (c *AdminController) IndexPage(f *fiber.Ctx) error {
	return f.Render("index", fiber.Map{
		"Title": "Головна",
	})
}

func (c *AdminController) MeasurementsPage(f *fiber.Ctx) error {
	var (
		ctx = context.Background()
	)

	measurements, err := c.ms.List(ctx)
	if err != nil {
		return f.Status(fiber.StatusInternalServerError).SendString("Failed to load measurements")
	}

	return f.Render("measurement", fiber.Map{
		"Title":        "Виміри",
		"Measurements": measurements,
	})
}
