package middleware

import (
	"context"

	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
)

func NewAuthMiddleware(authClient *auth.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionCookie := c.Cookies("session")
		if sessionCookie == "" {
			return c.Redirect("/admin/login")
		}

		token, err := authClient.VerifySessionCookieAndCheckRevoked(context.Background(), sessionCookie)
		if err != nil {
			return c.Redirect("/admin/login")
		}

		c.Locals("user", token)

		return c.Next()
	}
}
