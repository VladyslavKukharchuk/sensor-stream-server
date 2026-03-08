package routes

import (
	"github.com/gofiber/fiber/v2"
	"sensor-stream-server/internal/controller/auth"
)

func RegisterAuthRoutes(r fiber.Router, authController *auth.Controller) {
	r.Get("/login", authController.LoginPage)
	r.Post("/auth/session", authController.CreateSession)
	r.Get("/logout", authController.Logout)
}
