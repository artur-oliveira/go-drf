package database

import (
	"grf/bootstrap/grf"
	"grf/domain/auth"
	"log"
)

func RegisterMigrations(app *grf.App) {
	log.Println("starting database migrations")

	var allModels []interface{}

	allModels = append(allModels, auth.GetModels()...)

	if err := PerformMigration(app.DB, app.Config, allModels...); err != nil {
		log.Fatalf("Failed to perform automigrations: %v", err)
	}

	log.Println("database migrations complete.")
}
