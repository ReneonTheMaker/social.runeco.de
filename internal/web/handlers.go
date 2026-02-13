package web

import (
	"app/internal/model"
	"app/internal/store"
	"strconv"

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

func PostAuth(store *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		if username == "" || password == "" {
			return c.Status(fiber.StatusBadRequest).Render("loginform", fiber.Map{
				"Error": "Username and password are required",
			})
		}

		user, err := store.AuthenticateUser(username, password)
		if err.Error() == "record not found" {
			newUser := model.User{
				Username:     username,
				PasswordHash: password, // In production, hash the password before storing
			}
			err = store.CreateUser(&newUser)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).Render("loginform", fiber.Map{
					"Error": "Internal server error",
				})
			}
			user, _ = store.AuthenticateUser(username, password)
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).Render("loginform", fiber.Map{
				"Error": "Internal server error",
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "user_id",
			Value:    strconv.Itoa(int(user.ID)),
			HTTPOnly: true,
			SameSite: "Lax",
		})

		c.Set("HX-Redirect", "/feed")
		return c.SendStatus(fiber.StatusOK)
	}
}
