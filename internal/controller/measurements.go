package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type MeasurementRequest struct {
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Timestamp   time.Time `json:"timestamp"`
}

func AddMeasurements(c *fiber.Ctx) error {
	var m MeasurementRequest

	if err := c.BodyParser(&m); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid JSON",
		})
	}

	return c.JSON(fiber.Map{
		"message": "measurement received",
	})
}
