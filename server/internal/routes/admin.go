package routes

import (
	"github.com/gofiber/fiber/v2"

	"sensor-stream-server/internal/controller/admin"
)

func RegisterAdminRoutes(r fiber.Router, adminController *admin.Controller, auth fiber.Handler) {
	r.Get("/login", adminController.LoginPage)

	r.Get("/", auth, adminController.IndexPage)
	r.Get("/devices/:id", auth, adminController.DevicePage)
	r.Post("/devices/:id", auth, adminController.UpdateDevice)
}
