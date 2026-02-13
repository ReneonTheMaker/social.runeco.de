package web

// middleware to create id cookie that sets a uuid if not present
import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func IdCookieMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Cookies("id")
		if id == "" {
			newId := uuid.New().String()
			c.Cookie(&fiber.Cookie{
				Name:  "id",
				Value: newId,
			})
		}
		return c.Next()
	}
}

func RegisterMiddleware(app *fiber.App) {
	app.Use(IdCookieMiddleware())
}
