package devices

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"sensor-stream-server/internal/model"
)

type Service interface {
	Add(context.Context, string) (*model.Device, error)
}

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

type DeviceRequest struct {
	MAC string `json:"mac"`
}

func (c *Controller) Add(f *fiber.Ctx) error {
	var (
		ctx = context.Background()
	)

	var m DeviceRequest

	if err := f.BodyParser(&m); err != nil {
		writeErr := f.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		if writeErr != nil {
			return fmt.Errorf("failed to write error response: %w", writeErr)
		}

		return nil
	}

	device, err := c.service.Add(ctx, m.MAC)
	if err != nil {
		return fmt.Errorf("failed to add device: %w", err)
	}

	writeOK := f.JSON(fiber.Map{
		"id":  device.ID,
		"mac": device.MAC,
	})
	if writeOK != nil {
		return fmt.Errorf("failed to write success response: %w", writeOK)
	}

	return nil
}
