package routes

import (
	"github.com/gofiber/fiber/v2"

	"sensor-stream-server/internal/controller"
)

func RegisterAdminRoutes(app *fiber.App, adminController *controller.AdminController) {
	admin := app.Group("/admin")

	admin.Get("/", adminController.IndexPage)
	admin.Get("/measurements", adminController.MeasurementsPage)
}
