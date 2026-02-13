package web

import (
	"github.com/gofiber/fiber/v2"
)

// index page handler
func GetMain() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		// get number of hits from the store
		return c.Render("index", nil)
	}
}

func GetAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.Render("auth", nil)
	}
}
