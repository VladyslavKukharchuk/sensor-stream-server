package controller

import (
	"fmt"
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
		writeErr := c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		if writeErr != nil {
			return fmt.Errorf("failed to write error response: %w", writeErr)
		}

		return nil
	}

	writeOK := c.JSON(fiber.Map{"status": "ok"})
	if writeOK != nil {
		return fmt.Errorf("failed to write success response: %w", writeOK)
	}

	return nil
}
