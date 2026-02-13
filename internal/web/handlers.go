package web

import (
	"app/internal/store"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// index page handler
func GetMain(hitStore *store.HitsStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		// get number of hits from the store
		return c.Render("index", nil)
	}
}

func GetHits(hitStore *store.HitsStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		id := c.Cookies("id")
		if id == "" {
			return c.SendString("0")
		}
		hits, _ := hitStore.GetHits(id)
		return c.SendString(strconv.Itoa(hits))
	}
}

func PostHits(hitStore *store.HitsStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		time.Sleep(2 * time.Second) // Simulate a delay for fetching IP information
		id := c.Cookies("id")
		if id == "" {
			return c.SendString("0")
		}
		hits, _ := hitStore.GetHits(id)
		hits++
		hitStore.SetHits(id, hits)
		return c.SendString(strconv.Itoa(hits))
	}
}

func GetIpPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html")
	return c.Render("ipPage", nil)
}

func GetIp(c *fiber.Ctx) error {
	ip := c.IP()
	time.Sleep(2 * time.Second) // Simulate a delay for fetching IP information
	return c.SendString(ip)
}
