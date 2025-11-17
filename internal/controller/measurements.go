package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"sensor-stream-server/internal/model"
)

type Service interface {
	Add(ctx context.Context, m *model.Measurement) error
}

type MeasurementController struct {
	service Service
}

func NewMeasurementController(service Service) *MeasurementController {
	return &MeasurementController{service: service}
}

type MeasurementRequest struct {
	DeviceID    string    `json:"device_id"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Timestamp   time.Time `json:"timestamp"`
}

func (m MeasurementRequest) toMeasurementModel() *model.Measurement {
	return &model.Measurement{
		DeviceID:    m.DeviceID,
		Temperature: m.Temperature,
		Humidity:    m.Humidity,
		Timestamp:   m.Timestamp,
	}
}

func (c *MeasurementController) Add(f *fiber.Ctx) error {
	var (
		ctx = context.Background()
	)

	var m MeasurementRequest

	if err := f.BodyParser(&m); err != nil {
		writeErr := f.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		if writeErr != nil {
			return fmt.Errorf("failed to write error response: %w", writeErr)
		}

		return nil
	}

	err := c.service.Add(ctx, m.toMeasurementModel())
	if err != nil {
		return fmt.Errorf("failed to add measurement: %w", err)
	}

	writeOK := f.JSON(fiber.Map{"status": "ok"})
	if writeOK != nil {
		return fmt.Errorf("failed to write success response: %w", writeOK)
	}

	return nil
}
