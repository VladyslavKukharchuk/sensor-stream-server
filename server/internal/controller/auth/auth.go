package auth

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
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

func (c *Controller) DestroySession(f *fiber.Ctx) error {
	f.ClearCookie("session")
	return f.SendStatus(fiber.StatusOK)
}
