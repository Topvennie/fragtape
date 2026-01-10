// Package middlewares contains the server middlewares
package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
	"go.uber.org/zap"
)

func ProtectedRoute(c *fiber.Ctx) error {
	session, err := goth_fiber.SessionStore.Get(c)
	if err != nil {
		zap.S().Errorf("Failed to get session %w", err)
		return fiber.ErrInternalServerError
	}

	if session.Fresh() {
		return c.Redirect("/", fiber.StatusUnauthorized)
	}

	var userID any
	if userID = session.Get("userID"); userID == nil {
		return c.Redirect("/", fiber.StatusForbidden)
	}

	var steamUID any
	if steamUID = session.Get("steamUID"); steamUID == nil {
		return c.Redirect("/")
	}

	c.Locals("userID", userID)
	c.Locals("steamUID", steamUID)

	return c.Next()
}
