package web

import (
	"app/internal/store"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, hitsStore *store.HitsStore) {
	app.Get("/", GetMain(hitsStore))
	app.Get("/hits", GetHits(hitsStore))
	app.Post("/hits", PostHits(hitsStore))
	app.Get("/ip", GetIpPage)
	app.Get("/userinfo/ip", GetIp)
}
