package main

import (
	"log"

	"app/internal/config"
	"app/internal/web"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config.ini")
	if err != nil {
		log.Fatal(err)
	}
	app := web.NewApp()
	log.Fatal(app.FiberApp.Listen(":" + cfg.Web.Port))
}
