package routes

import (
	"github.com/gofiber/fiber/v2"

	"sensor-stream-server/internal/controller"
)

func RegisterMeasurementRoutes(
	app *fiber.App,
	measurementController *controller.MeasurementController,
	devicesController *controller.DevicesController,
) {
	api := app.Group("/api/v1")

	api.Post("/measurements", measurementController.Add)
	api.Post("/devices", devicesController.Add)
}
