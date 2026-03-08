package routes

import (
	"github.com/gofiber/fiber/v2"

	"sensor-stream-server/internal/controller/auth"
)

func RegisterAuthRoutes(r fiber.Router, authController *auth.Controller) {
	r.Post("/session", authController.CreateSession)
	r.Delete("/session", authController.DestroySession)
}
