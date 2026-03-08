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
			return c.Redirect("/login")
		}

		// Verify the ID token from the cookie
		token, err := authClient.VerifyIDToken(context.Background(), sessionCookie)
		if err != nil {
			return c.Redirect("/login")
		}

		// Store user info in context for controllers
		c.Locals("user", token)

		return c.Next()
	}
}
