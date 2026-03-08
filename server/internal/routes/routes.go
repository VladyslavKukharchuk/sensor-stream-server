package routes

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"

	"sensor-stream-server/internal/controller/admin"
	authController "sensor-stream-server/internal/controller/auth"
	"sensor-stream-server/internal/controller/devices"
	"sensor-stream-server/internal/controller/measurements"
	"sensor-stream-server/internal/middleware"
)

func Setup(
	app *fiber.App,
	authClient *auth.Client,
	mc *measurements.Controller,
	dc *devices.Controller,
	ac *admin.Controller,
	auc *authController.Controller,
) {
	app.Get("/login", ac.LoginPage)

	authGroup := app.Group("/auth")
	RegisterAuthRoutes(authGroup, auc)

	RegisterMeasurementRoutes(app, mc, dc)

	adminGroup := app.Group("/admin", middleware.NewAuthMiddleware(authClient))
	RegisterAdminRoutes(adminGroup, ac)
}
