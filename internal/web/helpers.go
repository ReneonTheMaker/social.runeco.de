package web

import (
	"time"

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
