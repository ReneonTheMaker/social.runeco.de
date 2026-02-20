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
		if err != nil && err.Error() == "record not found" {
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

func GetFeed() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		posts := []model.Post{
			{ID: 1, Content: "Hello, world!"},
			{ID: 2, Content: "Welcome to the feed!"},
		}
		return c.Render("feed", fiber.Map{
			"Posts": posts,
		})
	}
}

func PostDeletePost(store *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// post id from url param
		postID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).Render("feed", fiber.Map{
				"Error": "Invalid post ID",
			})
		}
		err = store.DeletePost(uint(postID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).Render("feed", fiber.Map{
				"Error": "Internal server error",
			})
		}
		c.Set("HX-Redirect", "/feed")
		return c.SendStatus(fiber.StatusOK)
	}
}
