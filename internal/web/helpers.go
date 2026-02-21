package web

import (
	"time"

	"app/internal/store"

	"github.com/gofiber/fiber/v2"
)

func SetCookie(c *fiber.Ctx, name, value string, duration time.Duration) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    value,
		HTTPOnly: true,
		SameSite: "Lax",
		Expires:  time.Now().Add(365 * 24 * time.Hour), // 1 year
		Secure:   true,                                 // Uncomment if using HTTPS
	})
}

func ClearCookie(c *fiber.Ctx, name string) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    "",
		HTTPOnly: true,
		SameSite: "Lax",
		Expires:  time.Now().Add(-1 * time.Hour), // Expire in the past
		Secure:   true,                           // Uncomment if using HTTPS
	})
}

func CheckAuth(c *fiber.Ctx, store *store.Store) error {
	// Check if user_id cookie is present and valid
	if c.Cookies("session") == "" {
		return c.Redirect("/auth", fiber.StatusSeeOther)
	} else {
		sessionId := c.Cookies("session")
		if !Auth(sessionId, store) {
			ClearCookie(c, "session")
			return c.Redirect("/auth", fiber.StatusSeeOther)
		}
	}
	return nil
}
