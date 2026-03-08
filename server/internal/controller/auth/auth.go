package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Config struct {
	FirebaseApiKey     string
	FirebaseAuthDomain string
	FirebaseProjectId  string
}

type Controller struct {
	config Config
}

func NewController(config Config) *Controller {
	return &Controller{config: config}
}

func (c *Controller) LoginPage(f *fiber.Ctx) error {
	return f.Render("login", fiber.Map{
		"FirebaseApiKey":     c.config.FirebaseApiKey,
		"FirebaseAuthDomain": c.config.FirebaseAuthDomain,
		"FirebaseProjectId":  c.config.FirebaseProjectId,
	})
}

func (c *Controller) CreateSession(f *fiber.Ctx) error {
	type Request struct {
		IDToken string `json:"idToken"`
	}

	var req Request
	if err := f.BodyParser(&req); err != nil {
		return f.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	f.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    req.IDToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	return f.SendStatus(fiber.StatusOK)
}

func (c *Controller) Logout(f *fiber.Ctx) error {
	f.ClearCookie("session")
	return f.Redirect("/login")
}
