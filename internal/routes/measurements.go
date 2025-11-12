package routes

import (
	"github.com/gofiber/fiber/v2"

	"sensor-stream-server/internal/controller"
)

func RegisterSensorRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Post("/measurements", controller.AddMeasurements)
}
