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
	"sensor-stream-server/internal/middleware"
	"sensor-stream-server/internal/repository"
	"sensor-stream-server/internal/routes"
	"sensor-stream-server/internal/service"
)

func main() {
	_ = godotenv.Load()

	projectID := os.Getenv("FIRESTORE_PROJECT_ID")
	if projectID == "" {
		log.Fatal().Msg("FIRESTORE_PROJECT_ID is not set")
	}

	databaseID := os.Getenv("FIRESTORE_DATABASE_ID")
	if databaseID == "" {
		log.Fatal().Msg("FIRESTORE_DATABASE_ID is not set")
	}

	firebaseApiKey := os.Getenv("FIREBASE_API_KEY")
	if firebaseApiKey == "" {
		log.Fatal().Msg("FIREBASE_API_KEY is not set")
	}

	firebaseAuthDomain := os.Getenv("FIREBASE_AUTH_DOMAIN")
	if firebaseAuthDomain == "" {
		log.Fatal().Msg("FIREBASE_AUTH_DOMAIN is not set")
	}

	engine := html.New("./internal/views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./public")
	app.Use(logger.New())
	app.Use(cors.New())

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
	measurementService := service.NewMeasurementService(measurementRepo)
	measurementController := measurements.NewController(measurementService)

	devicesRepo := repository.NewDevicesRepository(firestoreClient)
	devicesService := service.NewDevicesService(devicesRepo)
	devicesController := devices.NewController(devicesService)

	adminController := admin.NewController(measurementService, devicesService)

	authConfig := auth.Config{
		FirebaseApiKey:     firebaseApiKey,
		FirebaseAuthDomain: firebaseAuthDomain,
		FirebaseProjectId:  projectID,
	}
	authController := auth.NewController(authConfig)

	// Routes
	app.Get("/login", authController.LoginPage)
	app.Post("/auth/session", authController.CreateSession)
	app.Get("/logout", authController.Logout)

	routes.RegisterMeasurementRoutes(app, measurementController, devicesController)

	adminGroup := app.Group("/admin", middleware.NewAuthMiddleware(authClient))
	routes.RegisterAdminRoutes(adminGroup, adminController)

	if err := app.Listen(":8080"); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
