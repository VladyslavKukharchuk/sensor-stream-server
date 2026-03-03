package routes

import (
	"github.com/gofiber/fiber/v2"

	"sensor-stream-server/internal/controller/admin"
)

func RegisterAdminRoutes(app *fiber.App, adminController *admin.Controller) {
	r := app.Group("/admin")

	r.Get("/", adminController.IndexPage)
	r.Get("/measurements", adminController.MeasurementsPage)
	r.Get("/devices", adminController.DevicesPage)
	r.Get("/devices/:id", adminController.DevicePage)
}
