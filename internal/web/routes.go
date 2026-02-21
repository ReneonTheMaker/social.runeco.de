package web

import (
	"app/internal/store"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, store *store.Store) {
	app.Get("/", GetMain())
	app.Get("/auth", GetAuth(store))
	app.Post("/auth", PostAuth(store))
	app.Get("/feed", GetFeed(store))

	// post routes
	app.Get("/post/:id/reply-count", GetPostNumberOfReplies(store))
	app.Get("/post/:id", GetPost(store))
	app.Post("/post/:id/delete", PostDeletePost(store))
	app.Post("/post/:id/reply", PostReply(store))
	app.Post("/post/create", PostCreatePost(store))
}
