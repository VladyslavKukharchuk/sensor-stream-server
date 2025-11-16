package routes

import (
	"github.com/gofiber/fiber/v2"

	"sensor-stream-server/internal/controller"
)

func RegisterMeasurementRoutes(app *fiber.App, measurementController *controller.MeasurementController) {
	api := app.Group("/api/v1")

	api.Post("/measurements", measurementController.Add)
}
