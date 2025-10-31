package main

import (
	"grf/bootstrap"
	"log"
)

func main() {
	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(app.FiberApp.Listen(":" + app.Config.DBPort))
}
