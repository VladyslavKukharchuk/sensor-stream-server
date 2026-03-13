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
	"sensor-stream-server/internal/controller/auth"
	"sensor-stream-server/internal/controller/devices"
	"sensor-stream-server/internal/controller/measurements"
	"sensor-stream-server/internal/db"
	"sensor-stream-server/internal/repository"
	"sensor-stream-server/internal/routes"
	"sensor-stream-server/internal/service"
)

func main() {
	_ = godotenv.Load()

	projectID := getEnvOrFatal("FIRESTORE_PROJECT_ID")
	databaseID := getEnvOrFatal("FIRESTORE_DATABASE_ID")
	firebaseApiKey := getEnvOrFatal("FIREBASE_API_KEY")
	firebaseAuthDomain := getEnvOrFatal("FIREBASE_AUTH_DOMAIN")

	engine := html.New("./internal/views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())
	app.Use(cors.New())
	app.Static("/", "./public")

	ctx := context.Background()

	firestoreClient, err := db.NewFirestore(ctx, projectID, databaseID)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create firestore client")
	}

	authClient, err := db.NewFirebaseAuth(ctx, projectID)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create firebase auth client")
	}

	measurementRepo := repository.NewMeasurementRepository(firestoreClient)
	devicesRepo := repository.NewDevicesRepository(firestoreClient)

	measurementService := service.NewMeasurementService(measurementRepo)
	devicesService := service.NewDevicesService(devicesRepo)

	mc := measurements.NewController(measurementService)
	dc := devices.NewController(devicesService)
	ac := admin.NewController(measurementService, devicesService, admin.Config{
		FirebaseApiKey:     firebaseApiKey,
		FirebaseAuthDomain: firebaseAuthDomain,
		FirebaseProjectId:  projectID,
	})
	auc := auth.NewController()

	routes.Setup(app, authClient, mc, dc, ac, auc)

	if err := app.Listen(":8080"); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

func getEnvOrFatal(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatal().Msgf("%s is not set", key)
	}

	return val
}
