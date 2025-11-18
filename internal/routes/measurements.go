package routes

import (
	"github.com/gofiber/fiber/v2"

	"sensor-stream-server/internal/controller/devices"
	"sensor-stream-server/internal/controller/measurements"
)

func RegisterMeasurementRoutes(
	app *fiber.App,
	measurementController *measurements.Controller,
	devicesController *devices.Controller,
) {
	r := app.Group("/api/v1")

	r.Post("/measurements", measurementController.Add)
	r.Post("/devices", devicesController.Add)
}
