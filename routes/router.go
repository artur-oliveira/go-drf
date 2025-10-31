package routes

import (
	"grf/bootstrap/grf"
)

func RegisterRoutes(app *grf.App) {
	apiV1 := app.FiberApp.Group("/v1")
	RegisterAuthRoutes(apiV1, app)
}
