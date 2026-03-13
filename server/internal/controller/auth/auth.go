//nolint:wrapcheck
package auth

import (
	"context"
	"fmt"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
)

const (
	// Firebase session cookies can last up to 14 days.
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

	// Create a Firebase session cookie
	expiresIn := time.Duration(sessionDurationDays) * hoursInDay * time.Hour

	sessionCookie, err := c.authClient.SessionCookie(context.Background(), req.IDToken, expiresIn)
	if err != nil {
		return f.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": fmt.Sprintf("failed to create session cookie: %v", err)})
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
