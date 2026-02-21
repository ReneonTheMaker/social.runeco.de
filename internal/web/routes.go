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
	app.Post("/post/:id/delete", PostDeletePost(store))
}
