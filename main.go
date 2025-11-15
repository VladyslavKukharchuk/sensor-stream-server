package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/rs/zerolog/log"

	"sensor-stream-server/internal/routes"
)

func main() {
	engine := html.New("./internal/views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Sensor Dashboard",
		})
	})

	routes.RegisterSensorRoutes(app)

	if err := app.Listen(":8080"); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
