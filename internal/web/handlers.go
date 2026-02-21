package web

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"app/internal/store"

	"github.com/gofiber/fiber/v2"
)

// index page handler
func GetMain() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		// get number of hits from the store
		if c.Cookies("session") != "" {
			return c.Render("index", fiber.Map{
				"Authenticated": true,
			})
		}
		return c.Render("index", nil)
	}
}

func GetAuth(store *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Cookies("session") != "" {
			sessionId := c.Cookies("session")
			if Auth(sessionId, store) {
				return c.Redirect("/feed", fiber.StatusSeeOther)
			}
		}
		c.Set("Content-Type", "text/html")
		return c.Render("auth", nil)
	}
}

func PostAuth(store *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		if username == "" || password == "" {
			return c.Status(fiber.StatusBadRequest).Render("auth", fiber.Map{
				"Error": "Username and password are required",
			})
		}

		user, err := LoginOrSignUp(username, password, store)
		if err != nil {
			log.Printf("Error during login/signup: %v", err)
			return c.Status(fiber.StatusInternalServerError).Render("auth", fiber.Map{
				"Error": "Internal server error",
			})
		}

		sessionId := store.SetUserSession(user.ID)
		SetCookie(c, "session", sessionId, 24*time.Hour)

		c.Set("HX-Redirect", "/feed")
		return c.SendStatus(fiber.StatusOK)
	}
}

func GetFeed(store *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		CheckAuth(c, store)

		posts, err := store.GetTopPosts(10)
		if err != nil {
			log.Printf("Error fetching posts: %v", err)
			return c.Status(fiber.StatusInternalServerError).Render("feed", fiber.Map{
				"Error": "Internal server error",
			})
		}

		for i := range posts {
			user, err := store.GetUserFromSession(c.Cookies("session"))
			if err != nil {
				log.Printf("Error getting user ID from session: %v", err)
				continue
			}
			posts[i].Deleteable = posts[i].UserID == user.ID || user.Mod
		}

		return c.Render("feed", fiber.Map{
			"Authenticated": true,
			"Posts":         posts,
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

func GetPostNumberOfReplies(store *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		postID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid post ID",
			})
		}
		count, err := store.GetNumberOfReplies(uint(postID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
		c.Set("Content-Type", "text/html")
		return c.SendString(fmt.Sprintf("Replies: %d", count))
	}
}

func GetPost(store *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		CheckAuth(c, store)

		postID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Render("post_page", fiber.Map{})
		}
		post, err := store.GetPostByID(uint(postID))
		if err != nil {
			return c.Render("post_page", fiber.Map{})
		}
		replies, err := store.GetReplyPosts(uint(postID))
		if err != nil {
			return c.Render("post_page", fiber.Map{})
		}

		return c.Render("post_page", fiber.Map{
			"Post":          post,
			"Authenticated": true,
			"Replies":       replies,
		})
	}
}

func PostReply(store *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		CheckAuth(c, store)

		postID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).Render("post_page", fiber.Map{
				"Error": "Invalid post ID",
			})
		}
		content := c.FormValue("content")
		if content == "" {
			return c.Status(fiber.StatusBadRequest).Render("post_page", fiber.Map{
				"Error": "Content cannot be empty",
			})
		}
		id, err := store.GetUserIDFromSession(c.Cookies("session"))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).Render("post_page", fiber.Map{
				"Error": "Internal server error",
			})
		}
		reply, err := store.CreateReply(id, uint(postID), content)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).Render("post_page", fiber.Map{
				"Error": "Internal server error",
			})
		}
		return c.Render("post-reply", fiber.Map{
			"User":      reply.User,
			"CreatedAt": reply.CreatedAt,
			"Content":   reply.Content,
		})
	}
}

func PostCreatePost(store *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		CheckAuth(c, store)

		content := c.FormValue("content")

		if content == "" {
			return c.Status(fiber.StatusBadRequest).Render("feed", fiber.Map{
				"Error": "Content cannot be empty",
			})
		}

		id, err := store.GetUserIDFromSession(c.Cookies("session"))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).Render("feed", fiber.Map{
				"Error": "Internal server error",
			})
		}

		post, err := store.CreatePost(id, content)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).Render("feed", fiber.Map{
				"Error": "Internal server error",
			})
		}

		c.Set("content-type", "text/html")
		return c.Render("post", fiber.Map{
			"User":      post.User,
			"CreatedAt": post.CreatedAt,
			"Content":   post.Content,
			"ID":        post.ID,
		})
	}
}

func GetLogout(store *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionId := c.Cookies("session")
		if sessionId != "" {
			store.EndSession(sessionId)
			ClearCookie(c, "session")
			return c.Redirect("/", fiber.StatusSeeOther)
		}
		return c.Redirect("/", fiber.StatusSeeOther)
	}
}
