package web

import (
	"app/internal/db"
	"app/internal/model"
	"app/internal/store"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

type App struct {
	FiberApp *fiber.App
	Store    *store.Store
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

	// Favicon
	app.Static("/favicon.ico", "./static/favicon.ico")

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Initialize Store
	database := db.New("app.db")
	err := database.AutoMigrate(
		&model.User{},
		&model.UserInfo{},
		&model.Post{},
	)
	if err != nil {
		panic(err)
	}

	// Middleware to set ID cookie
	RegisterMiddleware(app)

	// Store
	storeInstance := store.NewStore(database)

	// Register Routes - defined in routes.go
	RegisterRoutes(app, storeInstance)

	return &App{
		FiberApp: app,
		Store:    storeInstance,
	}
}
