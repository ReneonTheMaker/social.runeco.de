package web

import (
	"app/internal/store"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

type App struct {
	FiberApp    *fiber.App
	HitStore    *store.HitsStore
}

func NewApp() *App {

	// Template Engine golang html/template
	engine := html.New("./views", ".html")

	registerRenderFunctions(engine)

	// Fiber App
	app := fiber.New(
		fiber.Config{
			Views: engine,
		},
	)

	// Static files
	app.Static("/static", "./static")

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Create the hits store
	hitStore := store.NewHitsStore()

	// Middleware to set ID cookie
	RegisterMiddleware(app)

	// Register Routes - defined in routes.go
	RegisterRoutes(app, hitStore)
	
	return &App{
		FiberApp: app,
		HitStore: hitStore,
	}
}
