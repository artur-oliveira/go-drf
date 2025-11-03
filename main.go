package main

import (
	"grf/bootstrap"
	"grf/config"
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
