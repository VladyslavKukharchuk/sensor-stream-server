package main

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"sensor-stream-server/internal/controller/admin"
	"sensor-stream-server/internal/controller/devices"
	"sensor-stream-server/internal/controller/measurements"
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

	app.Static("/", "./public")

	app.Use(logger.New())
	app.Use(cors.New())

	_ = godotenv.Load()

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
	measurementController := measurements.NewController(measurementService)
	devicesRepo := repository.NewDevicesRepository(firestoreClient)
	devicesService := service.NewDevicesService(devicesRepo)
	devicesController := devices.NewController(devicesService)
	adminController := admin.NewController(measurementService, devicesService)

	routes.RegisterMeasurementRoutes(app, measurementController, devicesController)
	routes.RegisterAdminRoutes(app, adminController)

	if err := app.Listen(":8080"); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
