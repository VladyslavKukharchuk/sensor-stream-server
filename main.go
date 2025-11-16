package main

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/rs/zerolog/log"

	"sensor-stream-server/internal/controller"
	"sensor-stream-server/internal/db"
	"sensor-stream-server/internal/repository"
	"sensor-stream-server/internal/routes"
	"sensor-stream-server/internal/service"
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

	firestoreProjectID := os.Getenv("FIRESTORE_PROJECT_ID")
	if firestoreProjectID == "" {
		log.Fatal().Msg("FIRESTORE_PROJECT_ID is not set")
	}

	firestoreDatabaseID := os.Getenv("FIRESTORE_DATABASE_ID")
	if firestoreDatabaseID == "" {
		log.Fatal().Msg("FIRESTORE_DATABASE_ID is not set")
	}

	ctx := context.Background()

	firestoreClient, err := db.NewFirestore(ctx, firestoreProjectID, firestoreDatabaseID)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create firestore client")
	}

	measurementRepo := repository.NewMeasurementRepository(firestoreClient)
	measurementService := service.NewMeasurementService(measurementRepo)
	measurementController := controller.NewMeasurementController(measurementService)

	routes.RegisterMeasurementRoutes(app, measurementController)

	if err := app.Listen(":8080"); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
