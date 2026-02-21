package web

import (
	"log"
	"strconv"

	"app/internal/model"
	"app/internal/store"

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

		user, err := store.GetUserByUsername(username)
		if err != nil {
			log.Printf("Error fetching user: %v", err)
		}
		if user == nil {
			// create user
			passwordHash, err := HashString(password)
			if err != nil {
				log.Printf("Error hashing password: %v", err)
				return c.Status(fiber.StatusInternalServerError).Render("auth", fiber.Map{
					"Error": "Internal server error",
				})
			}
			user = &model.User{
				Username:     username,
				PasswordHash: passwordHash,
			}
			err = store.CreateUser(user)
			if err != nil {
				log.Printf("Error creating user: %v", err)
				return c.Status(fiber.StatusInternalServerError).Render("auth", fiber.Map{
					"Error": "Internal server error",
				})
			}
		} else {
			// check password
			match, err := CheckPasswordHash(*user, password)
			if err != nil || !match {
				return c.Status(fiber.StatusUnauthorized).Render("auth", fiber.Map{
					"Error": "Invalid username or password",
				})
			}
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

func GetFeed(store *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		// Check if user_id cookie is present and valid
		if c.Cookies("user_id") == "" {
			return c.Redirect("/auth", fiber.StatusSeeOther)
		} else {
			id := c.Cookies("user_id")
			uId, err := strconv.Atoi(id)
			if err != nil {
				return c.Redirect("/auth", fiber.StatusSeeOther)
			}
			if !Auth(uint(uId), store) {
				return c.Redirect("/auth", fiber.StatusSeeOther)
			}
		}

		posts := []model.Post{
			{ID: 1, Content: "Hello, world!"},
			{ID: 2, Content: "Welcome to the feed!"},
			{ID: 3},
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
