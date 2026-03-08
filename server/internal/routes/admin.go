package routes

import (
	"github.com/gofiber/fiber/v2"

	"sensor-stream-server/internal/controller/admin"
)

func RegisterAdminRoutes(r fiber.Router, adminController *admin.Controller) {
	r.Get("/", adminController.IndexPage)
	r.Get("/measurements", adminController.MeasurementsPage)
	r.Get("/devices", adminController.DevicesPage)
	r.Get("/devices/:id", adminController.DevicePage)
	r.Post("/devices/:id", adminController.UpdateDevice)
}
