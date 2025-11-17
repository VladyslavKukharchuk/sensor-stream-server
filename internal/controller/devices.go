package controller

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"sensor-stream-server/internal/model"
)

type DevicesService interface {
	Add(context.Context, string) (*model.Device, error)
}

type DevicesController struct {
	service DevicesService
}

func NewDevicesController(service DevicesService) *DevicesController {
	return &DevicesController{service: service}
}

type DeviceRequest struct {
	MAC string `json:"mac"`
}

func (c *DevicesController) Add(f *fiber.Ctx) error {
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
