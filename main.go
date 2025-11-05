package main

import (
	"grf/core/bootstrap"
	"grf/core/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("./", "app")
	if err != nil {
		log.Fatal(err)
	}
	app, err := bootstrap.NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(app.FiberApp.Listen(":" + app.Config.ServerPort))
}
