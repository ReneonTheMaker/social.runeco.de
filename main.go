package main

import (
	"app/internal/config"
	"app/internal/web"
	"log"
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
