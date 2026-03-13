//nolint:wrapcheck
package auth

import (
	"fmt"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

const (
	sessionDurationDays = 14
	hoursInDay          = 24
)

type Controller struct {
	authClient *auth.Client
}

func NewController(authClient *auth.Client) *Controller {
	return &Controller{
		authClient: authClient,
	}
}

func (c *Controller) CreateSession(f *fiber.Ctx) error {
	type Request struct {
		IDToken string `json:"idToken"`
	}

	var req Request
	if err := f.BodyParser(&req); err != nil {
		return f.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	expiresIn := time.Duration(sessionDurationDays) * hoursInDay * time.Hour

	sessionCookie, err := c.authClient.SessionCookie(f.Context(), req.IDToken, expiresIn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create Firebase session cookie")

		return f.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": fmt.Sprintf("failed to create session: %v", err)})
	}

	f.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    sessionCookie,
		Expires:  time.Now().Add(expiresIn),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
	})

	return f.SendStatus(fiber.StatusOK)
}

func (c *Controller) DestroySession(f *fiber.Ctx) error {
	f.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    "",
		Expires:  time.Now().Add(-hoursInDay * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
	})

	return f.SendStatus(fiber.StatusOK)
}
